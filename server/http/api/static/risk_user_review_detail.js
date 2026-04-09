// 常量定义，保持与 list.js 一致
const DECISION_MANUAL_REVIEW = 1;
const decisionMap = {
  0:  { label: 'Unrecognized',         tagClass: 'tag-info' },
  1:  { label: 'Manual Review',        tagClass: 'tag-warning' },
  2:  { label: 'Block',                tagClass: 'tag-danger' },
  3:  { label: 'Watchlist',            tagClass: 'tag-info' },
  4:  { label: 'Allow',                tagClass: 'tag-success' },
  10: { label: 'Resolved - Block',     tagClass: 'tag-danger' },
  11: { label: 'Resolved - Whitelist', tagClass: 'tag-success' },
  12: { label: 'Resolved - Watchlist', tagClass: 'tag-info' },
};

const state = {
  record: null,
  updateLoading: false,
};

// 获取所有 DOM 元素
const els = {
  get: (id) => document.getElementById(id),
  init() {
    this.notFoundCard = this.get('notFoundCard');
    this.detailCard = this.get('detailCard');
    this.backBtnTop = this.get('backBtnTop');
    this.backBtnBottom = this.get('backBtnBottom');
    this.updateBtn = this.get('updateBtn');
    this.confirmDialogBtn = this.get('confirmDialogBtn');
    this.cancelDialogBtn = this.get('cancelDialogBtn');
    this.dialogMask = this.get('dialogMask');
    this.newDecisionSelect = this.get('newDecisionSelect');
    this.toastWrap = this.get('toastWrap');
    this.loadingMask = this.get('loadingMask');
  }
};

// 工具函数
function showToast(message, type = 'success') {
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.textContent = message;
  els.toastWrap.appendChild(toast);
  setTimeout(() => toast.remove(), 2500);
}

function formatTime(ts) {
  return ts ? new Date(Number(ts) * 1000).toLocaleString() : '-';
}

function renderRecord() {
  const r = state.record;
  if (!r) {
    els.notFoundCard.style.display = 'block';
    els.detailCard.style.display = 'none';
    return;
  }

  els.get('fieldId').textContent = r.id || '-';
  els.get('fieldUserId').textContent = r.user_id || '-';
  els.get('fieldCreateTime').textContent = formatTime(r.create_time);
  
  const config = decisionMap[Number(r.decision)] || { label: r.decision, tagClass: 'tag-info' };
  els.get('fieldDecision').innerHTML = `<span class="tag ${config.tagClass}">${config.label}</span>`;
  
  els.get('fieldDecisionSource').textContent = r.decision_source || '-';
  els.get('fieldRiskScore').textContent = r.risk_score ?? '-';
  els.get('fieldRiskLevel').textContent = r.risk_level || '-';
  els.get('fieldRuleScore').textContent = r.rule_score ?? '-';
  els.get('fieldFraudProbability').textContent = r.fraud_probability ?? '-';
  els.get('fieldAnalystSummary').textContent = r.analyst_summary || '-';
  els.get('fieldRules').textContent = r.rules || '-';

  els.updateBtn.style.display = Number(r.decision) === DECISION_MANUAL_REVIEW ? 'inline-block' : 'none';
  els.detailCard.style.display = 'block';
  els.notFoundCard.style.display = 'none';
}

// 交互逻辑
async function submitUpdate() {
  const newVal = els.newDecisionSelect.value;
  if (!newVal) return showToast('Please select a decision', 'warning');

  state.updateLoading = true;
  els.confirmDialogBtn.disabled = true;
  els.loadingMask.classList.add('show');

  try {
    const resp = await fetch(`/admin-ms/v1/merchant/risk-user-reviews/${state.record.id}/decision`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        id: state.record.id,
        user_id: state.record.user_id,
        decision: Number(newVal)
      })
    });

    const json = await resp.json();
    if (json.err_msg) {
      showToast(json.err_msg, 'error');
    } else {
      showToast('Updated successfully');
      state.record.decision = Number(newVal);
      sessionStorage.setItem('rur_detail_' + state.record.id, JSON.stringify(state.record));
      renderRecord();
      closeDialog();
    }
  } catch (e) {
    showToast('Request failed: ' + e.message, 'error');
  } finally {
    state.updateLoading = false;
    els.confirmDialogBtn.disabled = false;
    els.loadingMask.classList.remove('show');
  }
}

function openDialog() { els.dialogMask.classList.add('show'); }
function closeDialog() { els.dialogMask.classList.remove('show'); }
function goBack() { window.location.href = '/admin-ms/v1/merchant/risk-user-reviews/page'; }

// 初始化入口
document.addEventListener('DOMContentLoaded', () => {
  els.init();

  // 绑定事件
  els.backBtnTop.onclick = goBack;
  els.backBtnBottom.onclick = goBack;
  els.updateBtn.onclick = openDialog;
  els.cancelDialogBtn.onclick = closeDialog;
  els.confirmDialogBtn.onclick = submitUpdate;
  els.dialogMask.onclick = (e) => { if(e.target === els.dialogMask) closeDialog(); };

  // 获取数据
  const id = window.location.pathname.split('/').pop();
  const stored = sessionStorage.getItem('rur_detail_' + id);
  if (stored) {
    state.record = JSON.parse(stored);
  }
  renderRecord();
});