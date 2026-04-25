// ---- Constants ----
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

const featureNameMap = {
  order_count_last_1h:    'Orders in Last 1 Hour',
  order_count_last_24h:   'Orders in Last 24 Hours',
  unique_ip_count:        'Unique IP Count',
  account_age_days:       'Account Age',
  avg_order_amount_today: 'Avg Order Amount Today',
  receive_address_count:  'Receive Address Count',
};

const state = { record: null, updateLoading: false };

// ---- DOM element registry ----
const els = {
  get: (id) => document.getElementById(id),
  init() {
    this.notFoundCard       = this.get('notFoundCard');
    this.sectionSummary     = this.get('sectionSummary');
    this.sectionRisk        = this.get('sectionRisk');
    this.sectionAI          = this.get('sectionAI');
    this.sectionRules       = this.get('sectionRules');
    this.sectionRaw         = this.get('sectionRaw');
    this.actionBar          = this.get('actionBar');
    this.backBtnTop         = this.get('backBtnTop');
    this.backBtnBottom      = this.get('backBtnBottom');
    this.updateBtn          = this.get('updateBtn');
    this.confirmDialogBtn   = this.get('confirmDialogBtn');
    this.cancelDialogBtn    = this.get('cancelDialogBtn');
    this.dialogMask         = this.get('dialogMask');
    this.newDecisionSelect  = this.get('newDecisionSelect');
    this.toastWrap          = this.get('toastWrap');
    this.loadingMask        = this.get('loadingMask');
  }
};

// ---- Utility functions ----
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

function formatScore(value) {
  if (value === null || value === undefined || value === '') return 'N/A';
  const n = parseFloat(value);
  return isNaN(n) ? 'N/A' : n.toFixed(2);
}

function escHtml(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function formatSnakeCaseLabel(value) {
  if (!value) return String(value);
  if (featureNameMap[value]) return featureNameMap[value];
  return value.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
}

function safeParseJson(value) {
  if (value === null || value === undefined) return null;
  if (typeof value === 'object') return value;
  if (typeof value === 'string') {
    try { return JSON.parse(value); } catch (e) { return null; }
  }
  return null;
}

function normalizeRulesHit(value) {
  if (!value) return [];
  if (Array.isArray(value)) return value;
  const parsed = safeParseJson(value);
  if (Array.isArray(parsed)) return parsed;
  if (typeof value === 'string') {
    return value.split(',').map((s) => s.trim()).filter(Boolean);
  }
  return [];
}

function normalizeTopContributors(value) {
  if (!value) return null;
  const parsed = safeParseJson(value);
  if (!Array.isArray(parsed) || parsed.length === 0) return null;
  const items = parsed.map((item) => {
    const feature = item.feature || item.name || item.field || '';
    const contribution =
      item.contribution !== undefined ? item.contribution :
      item.score       !== undefined ? item.score :
      item.value       !== undefined ? item.value :
      item.weight      !== undefined ? item.weight :
      item.impact      !== undefined ? item.impact : null;
    return { feature, contribution, raw: item };
  });
  items.sort((a, b) =>
    Math.abs(parseFloat(b.contribution) || 0) - Math.abs(parseFloat(a.contribution) || 0)
  );
  return items;
}

function getRiskTagClass(riskLevel) {
  const lvl = (riskLevel || '').toLowerCase();
  if (lvl === 'low')                       return 'tag-risk-low';
  if (lvl === 'medium' || lvl === 'med')   return 'tag-risk-med';
  if (lvl === 'high')                      return 'tag-risk-high';
  return 'tag-info';
}

// ---- Render helpers ----
function renderDecisionTag(decision) {
  const cfg = decisionMap[Number(decision)] || { label: String(decision), tagClass: 'tag-info' };
  return `<span class="tag ${cfg.tagClass}">${escHtml(cfg.label)}</span>`;
}

function renderMlContributors(rawValue) {
  const contributors = normalizeTopContributors(rawValue);
  if (!contributors || contributors.length === 0) {
    return '<span class="no-data">No ML contribution details available.</span>';
  }
  let rows = '';
  for (const c of contributors) {
    const val = parseFloat(c.contribution);
    const isPos = !isNaN(val) && val > 0;
    const isNeg = !isNaN(val) && val < 0;
    const formatted = isNaN(val) ? 'N/A' : (isPos ? '+' : '') + val.toFixed(2);
    const cls = isPos ? 'contrib-pos' : (isNeg ? 'contrib-neg' : '');
    const ratio = c.raw && c.raw.ratio ? ` (${escHtml(String(c.raw.ratio))})` : '';
    let dirText;
    if (c.raw && c.raw.direction && typeof c.raw.direction === 'string') {
      dirText = escHtml(c.raw.direction);
    } else {
      dirText = isPos ? 'Increases Risk' : (isNeg ? 'Reduces Risk' : '-');
    }
    rows += `<tr>
      <td>${escHtml(formatSnakeCaseLabel(c.feature))}</td>
      <td class="${cls}">${escHtml(formatted)}${ratio}</td>
      <td>${dirText}</td>
    </tr>`;
  }
  return `<table class="contrib-table">
    <thead><tr><th>Feature</th><th>Contribution</th><th>Impact</th></tr></thead>
    <tbody>${rows}</tbody>
  </table>`;
}

function renderRulesTags(rulesValue) {
  const rules = normalizeRulesHit(rulesValue);
  if (rules.length === 0) return '<span class="no-data">No rules triggered.</span>';
  return '<div class="rules-tags">' +
    rules.map((r) => `<span class="rule-tag">${escHtml(formatSnakeCaseLabel(r))}</span>`).join('') +
    '</div>';
}

// ---- Main render ----
function renderRecord() {
  const r = state.record;
  if (!r) {
    els.notFoundCard.style.display = 'block';
    return;
  }

  // Section 1 – Review Summary
  els.get('fieldId').textContent        = r.id != null ? r.id : 'N/A';
  els.get('fieldUserId').textContent    = r.user_id != null ? r.user_id : 'N/A';
  els.get('fieldCreateTime').textContent = formatTime(r.create_time);
  els.get('fieldDecision').innerHTML    = renderDecisionTag(r.decision);
  els.get('fieldDecisionSource').textContent = r.decision_source || 'N/A';
  els.sectionSummary.style.display = 'block';

  // Section 2 – Risk Assessment
  const riskLvl = r.risk_level || '';
  els.get('fieldRiskLevel').innerHTML = riskLvl
    ? `<span class="tag ${getRiskTagClass(riskLvl)}">${escHtml(riskLvl.charAt(0).toUpperCase() + riskLvl.slice(1))}</span>`
    : '<span class="no-data">N/A</span>';
  els.get('fieldConfidence').textContent       = r.confidence || 'N/A';
  els.get('fieldRiskScore').textContent        = formatScore(r.risk_score);
  els.get('fieldRuleScore').textContent        = formatScore(r.rule_score);
  els.get('fieldFraudProbability').textContent = formatScore(r.fraud_probability);
  els.sectionRisk.style.display = 'block';

  // Section 3 – AI Explanation
  els.get('fieldAnalystSummary').textContent = r.analyst_summary || 'N/A';
  const mlRaw = r.ml_model_top_contributor
    || r.ml_top_contributors
    || r.top_contributors
    || r.top_contributions
    || null;
  els.get('fieldMlContributors').innerHTML = renderMlContributors(mlRaw);
  els.sectionAI.style.display = 'block';

  // Section 4 – Triggered Rules
  els.get('fieldRules').innerHTML = renderRulesTags(r.rules);
  els.sectionRules.style.display = 'block';

  // Section 5 – Raw Details
  els.get('fieldRawJson').textContent = JSON.stringify(r, null, 2);
  els.sectionRaw.style.display = 'block';

  // Action bar
  els.updateBtn.style.display = Number(r.decision) === DECISION_MANUAL_REVIEW ? 'inline-block' : 'none';
  els.actionBar.style.display = 'block';
}

// ---- Dialog ----
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
function goBack()      { window.location.href = '/admin-ms/v1/merchant/risk-user-reviews/page'; }

// ---- Init ----
document.addEventListener('DOMContentLoaded', () => {
  els.init();

  els.backBtnTop.onclick     = goBack;
  els.backBtnBottom.onclick  = goBack;
  els.updateBtn.onclick      = openDialog;
  els.cancelDialogBtn.onclick = closeDialog;
  els.confirmDialogBtn.onclick = submitUpdate;
  els.dialogMask.onclick     = (e) => { if (e.target === els.dialogMask) closeDialog(); };

  const id = window.location.pathname.split('/').pop();
  const stored = sessionStorage.getItem('rur_detail_' + id);
  if (stored) {
    try { state.record = JSON.parse(stored); } catch (e) {}
  }
  renderRecord();
});
