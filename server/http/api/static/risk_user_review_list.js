const DECISION_MANUAL_REVIEW = 1;
const DECISION_BLOCK = 2;
const DECISION_WATCHLIST = 3;
const DECISION_ALLOW = 4;
const DECISION_UNRECOGNIZED = 0;
const RESOLVED_BLOCK = 10;
const RESOLVED_WHITELIST = 11;
const RESOLVED_WATCHLIST = 12;

const decisionOptions = [
  { label: 'Unrecognized', value: DECISION_UNRECOGNIZED },
  { label: 'Manual Review', value: DECISION_MANUAL_REVIEW },
  { label: 'Block', value: DECISION_BLOCK },
  { label: 'Watchlist', value: DECISION_WATCHLIST },
  { label: 'Allow', value: DECISION_ALLOW },
  { label: 'Resolved - Block', value: RESOLVED_BLOCK },
  { label: 'Resolved - Whitelist', value: RESOLVED_WHITELIST },
  { label: 'Resolved - Watchlist', value: RESOLVED_WATCHLIST },
];

const state = {
  filter: {
    user_id: '',
    decision: '',
    start_time_ms: '',
    end_time_ms: '',
  },
  loading: false,
  tableData: [],
  total: 0,
  page: 1,
  pageSize: 20,
  selectedRow: null,
  updateLoading: false,
};

const els = {
  userIdInput: document.getElementById('userIdInput'),
  decisionSelect: document.getElementById('decisionSelect'),
  startTimeInput: document.getElementById('startTimeInput'),
  endTimeInput: document.getElementById('endTimeInput'),
  searchBtn: document.getElementById('searchBtn'),
  resetBtn: document.getElementById('resetBtn'),
  tableBody: document.getElementById('tableBody'),
  pageInfo: document.getElementById('pageInfo'),
  pageSizeSelect: document.getElementById('pageSizeSelect'),
  prevBtn: document.getElementById('prevBtn'),
  nextBtn: document.getElementById('nextBtn'),
  pageBtn: document.getElementById('pageBtn'),
  loadingMask: document.getElementById('loadingMask'),
  dialogMask: document.getElementById('dialogMask'),
  dialogReviewId: document.getElementById('dialogReviewId'),
  dialogUserId: document.getElementById('dialogUserId'),
  dialogCurrentDecision: document.getElementById('dialogCurrentDecision'),
  newDecisionSelect: document.getElementById('newDecisionSelect'),
  cancelDialogBtn: document.getElementById('cancelDialogBtn'),
  confirmDialogBtn: document.getElementById('confirmDialogBtn'),
  toastWrap: document.getElementById('toastWrap'),
};

function decisionLabel(val) {
  const opt = decisionOptions.find((o) => o.value === Number(val));
  return opt ? opt.label : String(val);
}

function decisionTagClass(val) {
  switch (Number(val)) {
    case DECISION_MANUAL_REVIEW:
      return 'tag-warning';
    case DECISION_BLOCK:
      return 'tag-danger';
    case DECISION_ALLOW:
      return 'tag-success';
    case RESOLVED_BLOCK:
      return 'tag-danger';
    case RESOLVED_WHITELIST:
      return 'tag-success';
    default:
      return 'tag-info';
  }
}

function formatTime(ts) {
  if (ts === null || ts === undefined || ts === '') return '-';
  return new Date(Number(ts) * 1000).toLocaleString();
}

function datetimeLocalToUnix(value) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return Math.floor(date.getTime() / 1000);
}

function setLoading(flag) {
  state.loading = flag;
  els.loadingMask.classList.toggle('show', flag);
  els.searchBtn.disabled = flag;
  els.resetBtn.disabled = flag;
  els.prevBtn.disabled = flag;
  els.nextBtn.disabled = flag;
  els.pageSizeSelect.disabled = flag;
}

function showToast(message, type) {
  const toast = document.createElement('div');
  toast.className = 'toast toast-' + type;
  toast.textContent = message;
  els.toastWrap.appendChild(toast);
  window.setTimeout(() => {
    if (toast.parentNode) {
      toast.parentNode.removeChild(toast);
    }
  }, 2500);
}

function renderDecisionOptions() {
  for (let i = 0; i < decisionOptions.length; i += 1) {
    const opt = decisionOptions[i];
    const option = document.createElement('option');
    option.value = String(opt.value);
    option.textContent = opt.label;
    els.decisionSelect.appendChild(option);
  }
}

function clearChildren(node) {
  while (node.firstChild) {
    node.removeChild(node.firstChild);
  }
}

function createCell(text, className) {
  const td = document.createElement('td');
  if (className) {
    td.className = className;
  }
  td.textContent = text === null || text === undefined || text === '' ? '-' : String(text);
  return td;
}

function renderTable() {
  clearChildren(els.tableBody);

  if (!state.tableData.length) {
    const tr = document.createElement('tr');
    const td = document.createElement('td');
    td.colSpan = 8;
    td.className = 'empty-text';
    td.textContent = 'No Data';
    tr.appendChild(td);
    els.tableBody.appendChild(tr);
    return;
  }

  for (let i = 0; i < state.tableData.length; i += 1) {
    const row = state.tableData[i];
    const tr = document.createElement('tr');

    tr.appendChild(createCell(row.id, 'nowrap-cell'));
    tr.appendChild(createCell(row.user_id, 'nowrap-cell'));
    tr.appendChild(createCell(formatTime(row.create_time), 'nowrap-cell'));
    tr.appendChild(createCell(row.risk_score, 'nowrap-cell'));
    tr.appendChild(createCell(row.risk_level, 'nowrap-cell'));

    const decisionTd = document.createElement('td');
    decisionTd.className = 'nowrap-cell';
    const tag = document.createElement('span');
    tag.className = 'tag ' + decisionTagClass(row.decision);
    tag.textContent = decisionLabel(row.decision);
    decisionTd.appendChild(tag);
    tr.appendChild(decisionTd);

    const sourceTd = document.createElement('td');
    sourceTd.className = 'nowrap-cell';
    sourceTd.title = row.decision_source || '-';
    sourceTd.textContent = row.decision_source || '-';
    tr.appendChild(sourceTd);

    const actionTd = document.createElement('td');
    actionTd.className = 'nowrap-cell';

    const actionWrap = document.createElement('div');
    actionWrap.className = 'table-actions';

    const detailBtn = document.createElement('button');
    detailBtn.className = 'btn';
    detailBtn.type = 'button';
    detailBtn.textContent = 'Detail';
    detailBtn.addEventListener('click', function () {
      viewDetail(row);
    });
    actionWrap.appendChild(detailBtn);

    if (Number(row.decision) === DECISION_MANUAL_REVIEW) {
      const updateBtn = document.createElement('button');
      updateBtn.className = 'btn btn-warning';
      updateBtn.type = 'button';
      updateBtn.textContent = 'Update';
      updateBtn.addEventListener('click', function () {
        openUpdateDialog(row);
      });
      actionWrap.appendChild(updateBtn);
    }

    actionTd.appendChild(actionWrap);
    tr.appendChild(actionTd);

    els.tableBody.appendChild(tr);
  }
}

function renderPagination() {
  const totalPages = Math.max(1, Math.ceil(state.total / state.pageSize));
  els.pageInfo.textContent = 'Total ' + state.total + ' · Page ' + state.page + ' / ' + totalPages;
  els.pageBtn.textContent = String(state.page);
  els.prevBtn.disabled = state.loading || state.page <= 1;
  els.nextBtn.disabled = state.loading || state.page >= totalPages;
}

function syncFilterFromInputs() {
  state.filter.user_id = els.userIdInput.value.trim();
  state.filter.decision = els.decisionSelect.value;
  state.filter.start_time_ms = els.startTimeInput.value;
  state.filter.end_time_ms = els.endTimeInput.value;
  state.pageSize = Number(els.pageSizeSelect.value);
}

async function fetchList(page) {
  if (typeof page === 'number') {
    state.page = page;
  }

  syncFilterFromInputs();
  setLoading(true);

  try {
    const params = new URLSearchParams();

    if (state.filter.user_id !== '') {
      params.set('user_id', state.filter.user_id);
    }
    if (state.filter.decision !== '') {
      params.set('decision', state.filter.decision);
    }

    const startTime = datetimeLocalToUnix(state.filter.start_time_ms);
    const endTime = datetimeLocalToUnix(state.filter.end_time_ms);

    if (startTime) {
      params.set('start_time', String(startTime));
    }
    if (endTime) {
      params.set('end_time', String(endTime));
    }

    params.set('page', String(state.page));
    params.set('page_size', String(state.pageSize));

    const resp = await fetch('/admin-ms/v1/merchant/risk-user-reviews?' + params.toString());
    const json = await resp.json();

    if (json.err_msg) {
      showToast(json.err_msg, 'error');
      state.tableData = [];
      state.total = 0;
    } else {
      state.tableData = json.data && json.data.list ? json.data.list : [];
      state.total = json.data && json.data.total ? json.data.total : 0;
    }
  } catch (e) {
    showToast('Request failed: ' + e.message, 'error');
    state.tableData = [];
    state.total = 0;
  } finally {
    setLoading(false);
    renderTable();
    renderPagination();
  }
}

function resetFilter() {
  els.userIdInput.value = '';
  els.decisionSelect.value = '';
  els.startTimeInput.value = '';
  els.endTimeInput.value = '';
  state.page = 1;
  fetchList(1);
}

function viewDetail(row) {
  sessionStorage.setItem('rur_detail_' + row.id, JSON.stringify(row));
  window.location.href = '/admin-ms/v1/merchant/risk-user-reviews/page/' + row.id;
}

function openUpdateDialog(row) {
  state.selectedRow = row;
  els.dialogReviewId.textContent = row.id ?? '-';
  els.dialogUserId.textContent = row.user_id ?? '-';

  clearChildren(els.dialogCurrentDecision);
  const tag = document.createElement('span');
  tag.className = 'tag ' + decisionTagClass(row.decision);
  tag.textContent = decisionLabel(row.decision);
  els.dialogCurrentDecision.appendChild(tag);

  els.newDecisionSelect.value = '';
  els.dialogMask.classList.add('show');
}

function closeUpdateDialog() {
  state.selectedRow = null;
  els.newDecisionSelect.value = '';
  els.dialogMask.classList.remove('show');
}

async function submitUpdateDecision() {
  const newDecision = els.newDecisionSelect.value;
  if (!newDecision) {
    showToast('Please select a resolved decision', 'warning');
    return;
  }
  if (!state.selectedRow) {
    return;
  }

  state.updateLoading = true;
  els.confirmDialogBtn.disabled = true;

  try {
    const body = {
      id: state.selectedRow.id,
      user_id: state.selectedRow.user_id,
      decision: Number(newDecision),
    };

    const resp = await fetch(
      '/admin-ms/v1/merchant/risk-user-reviews/' + state.selectedRow.id + '/decision',
      {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
      }
    );

    const json = await resp.json();

    if (json.err_msg) {
      showToast(json.err_msg, 'error');
    } else {
      showToast('Decision updated successfully', 'success');
      closeUpdateDialog();
      fetchList(state.page);
    }
  } catch (e) {
    showToast('Request failed: ' + e.message, 'error');
  } finally {
    state.updateLoading = false;
    els.confirmDialogBtn.disabled = false;
  }
}

els.searchBtn.addEventListener('click', function () {
  state.page = 1;
  fetchList(1);
});

els.resetBtn.addEventListener('click', resetFilter);

els.pageSizeSelect.addEventListener('change', function () {
  state.page = 1;
  fetchList(1);
});

els.prevBtn.addEventListener('click', function () {
  if (state.page > 1) {
    fetchList(state.page - 1);
  }
});

els.nextBtn.addEventListener('click', function () {
  const totalPages = Math.max(1, Math.ceil(state.total / state.pageSize));
  if (state.page < totalPages) {
    fetchList(state.page + 1);
  }
});

els.cancelDialogBtn.addEventListener('click', closeUpdateDialog);
els.confirmDialogBtn.addEventListener('click', submitUpdateDecision);

els.dialogMask.addEventListener('click', function (e) {
  if (e.target === els.dialogMask) {
    closeUpdateDialog();
  }
});

renderDecisionOptions();
renderPagination();
fetchList(1);