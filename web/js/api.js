const API = {
  async get(url) {
    const res = await fetch(url);
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    return res.json();
  },
  async post(url, data) {
    const res = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    });
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    return res.json();
  },
  getWeapons(params = {}) {
    const q = new URLSearchParams(params).toString();
    return this.get('/api/weapons' + (q ? '?' + q : ''));
  },
  getTop10() {
    return this.get('/api/weapons/top10');
  },
  getWeapon(id) {
    return this.get(`/api/weapons/${id}`);
  },
  compareWeapons(ids) {
    return this.post('/api/weapons/compare', { ids });
  },
  getAttachments(params = {}) {
    const q = new URLSearchParams(params).toString();
    return this.get('/api/attachments' + (q ? '?' + q : ''));
  },
  getModCodes(params = {}) {
    const q = new URLSearchParams(params).toString();
    return this.get('/api/mod-codes' + (q ? '?' + q : ''));
  },
  recommend(data) {
    return this.post('/api/recommend', data);
  },
  getKnowledge() {
    return this.get('/api/knowledge');
  }
};
