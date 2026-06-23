const SmartMod = {
  data() {
    return {
      weapons: [],
      selectedWeapon: null,
      style: 'balanced',
      budget: 500000,
      attrs: {
        effective_range: 0,
        vertical_recoil: 0,
        horiz_recoil: 0,
        handling_speed: 0,
        ads_stability: 0,
        hip_fire_acc: 0,
        muzzle_velocity: 0,
        sound_range: 0
      },
      results: [],
      loading: false,
      searched: false,
      copiedId: null,
      error: ''
    };
  },
  async created() {
    try {
      this.weapons = await API.getWeapons();
    } catch(e) { console.error(e); }
  },
  computed: {
    styles() {
      return [
        { key: 'balanced', label: '均衡' },
        { key: 'close', label: '近战突破' },
        { key: 'mid', label: '中距交战' },
        { key: 'long', label: '远程压制' },
        { key: 'stealth', label: '隐蔽作战' }
      ];
    },
    attrLabels() {
      return [
        { key: 'effective_range', label: '优势射程' },
        { key: 'vertical_recoil', label: '垂直后坐力' },
        { key: 'horiz_recoil', label: '水平后坐力' },
        { key: 'handling_speed', label: '操控速度' },
        { key: 'ads_stability', label: '据枪稳定性' },
        { key: 'hip_fire_acc', label: '腰射精度' },
        { key: 'muzzle_velocity', label: '枪口初速' },
        { key: 'sound_range', label: '消音能力' }
      ];
    }
  },
  methods: {
    selectWeapon(e) {
      const id = parseInt(e.target.value);
      this.selectedWeapon = this.weapons.find(w => w.id === id) || null;
      if (this.selectedWeapon) {
        this.attrs.effective_range = this.selectedWeapon.effective_range;
        this.attrs.vertical_recoil = this.selectedWeapon.vertical_recoil;
        this.attrs.horiz_recoil = this.selectedWeapon.horiz_recoil;
        this.attrs.handling_speed = this.selectedWeapon.handling_speed;
        this.attrs.ads_stability = this.selectedWeapon.ads_stability;
        this.attrs.hip_fire_acc = this.selectedWeapon.hip_fire_acc;
        this.attrs.muzzle_velocity = this.selectedWeapon.muzzle_velocity;
        this.attrs.sound_range = this.selectedWeapon.sound_range;
      }
      this.searched = false;
      this.results = [];
      this.error = '';
    },
    setStyle(s) {
      this.style = s;
      if (!this.selectedWeapon) return;
      const dir = {
        close:    {effective_range:-15, vertical_recoil:5, horiz_recoil:5, handling_speed:20, ads_stability:-5, hip_fire_acc:20, muzzle_velocity:-10, sound_range:-5},
        mid:      {effective_range:10, vertical_recoil:10, horiz_recoil:10, handling_speed:0, ads_stability:10, hip_fire_acc:0, muzzle_velocity:5, sound_range:0},
        long:     {effective_range:20, vertical_recoil:8, horiz_recoil:8, handling_speed:-15, ads_stability:15, hip_fire_acc:-10, muzzle_velocity:20, sound_range:0},
        balanced: {effective_range:5, vertical_recoil:5, horiz_recoil:5, handling_speed:5, ads_stability:5, hip_fire_acc:5, muzzle_velocity:5, sound_range:5},
        stealth:  {effective_range:0, vertical_recoil:3, horiz_recoil:3, handling_speed:0, ads_stability:5, hip_fire_acc:0, muzzle_velocity:-5, sound_range:25}
      };
      const d = dir[s] || dir.balanced;
      const w = this.selectedWeapon;
      const clamp = v => Math.max(0, Math.min(100, Math.round(v)));
      this.attrs.effective_range = clamp(w.effective_range + d.effective_range);
      this.attrs.vertical_recoil = clamp(w.vertical_recoil + d.vertical_recoil);
      this.attrs.horiz_recoil = clamp(w.horiz_recoil + d.horiz_recoil);
      this.attrs.handling_speed = clamp(w.handling_speed + d.handling_speed);
      this.attrs.ads_stability = clamp(w.ads_stability + d.ads_stability);
      this.attrs.hip_fire_acc = clamp(w.hip_fire_acc + d.hip_fire_acc);
      this.attrs.muzzle_velocity = clamp(w.muzzle_velocity + d.muzzle_velocity);
      this.attrs.sound_range = clamp(w.sound_range + d.sound_range);
    },
    async doRecommend() {
      if (!this.selectedWeapon) return;
      this.loading = true;
      this.error = '';
      try {
        const data = await API.recommend({
          weapon_id: this.selectedWeapon.id,
          budget: this.budget,
          style: this.style,
          ...this.attrs
        });
        this.results = Array.isArray(data) ? data : [];
        this.searched = true;
      } catch(e) {
        console.error(e);
        this.error = '请求失败，请检查服务是否正常运行';
        this.results = [];
        this.searched = true;
      }
      this.loading = false;
    },
    formatPrice(p) {
      if (p >= 10000) return (p / 10000).toFixed(1) + '万';
      return p.toLocaleString();
    },
    async copyCode(code, idx) {
      try {
        await navigator.clipboard.writeText(code);
      } catch(e) {
        const ta = document.createElement('textarea');
        ta.value = code;
        document.body.appendChild(ta);
        ta.select();
        document.execCommand('copy');
        document.body.removeChild(ta);
      }
      this.copiedId = idx;
      setTimeout(() => { this.copiedId = null; }, 2000);
    },
    resultRadar(r) {
      const mc = r.mod_code;
      const ds = [{
        name: '推荐方案',
        values: [mc.effective_range, mc.vertical_recoil, mc.horiz_recoil,
                 mc.handling_speed, mc.ads_stability, mc.hip_fire_acc,
                 mc.muzzle_velocity, mc.sound_range]
      }];
      if (this.selectedWeapon) {
        ds.push({
          name: '裸枪基础',
          values: [this.selectedWeapon.effective_range, this.selectedWeapon.vertical_recoil,
                   this.selectedWeapon.horiz_recoil, this.selectedWeapon.handling_speed,
                   this.selectedWeapon.ads_stability, this.selectedWeapon.hip_fire_acc,
                   this.selectedWeapon.muzzle_velocity, this.selectedWeapon.sound_range]
        });
      }
      return ds;
    },
    attrDiff(r, key) {
      if (!this.selectedWeapon) return 0;
      return Math.round((r.mod_code[key] - this.selectedWeapon[key]) * 10) / 10;
    },
    diffClass(v) {
      if (v > 0) return 'positive';
      if (v < 0) return 'negative';
      return '';
    },
    diffStr(v) {
      if (v > 0) return '+' + v;
      if (v < 0) return '' + v;
      return '0';
    },
    isDynamic(r) {
      return r.reason && r.reason.startsWith('[智能配件组合]');
    },
    countParts(r) {
      if (!r.mod_code.parts) return 0;
      return r.mod_code.parts.split('+').length;
    },
    tierBorder(grade) {
      const colors = {'满改':'rgba(255,215,0,0.25)','半改':'rgba(168,85,247,0.25)','丐版':'rgba(107,114,128,0.2)'};
      return {borderLeftWidth:'3px', borderLeftColor: colors[grade]||'transparent'};
    },
    scoreColor(s) {
      if (s >= 85) return '#32cd32';
      if (s >= 70) return '#fbbf24';
      return '#ef4444';
    }
  },
  template: `
    <div>
      <div class="section-header">
        <h2><span class="icon">&#9889;</span> 一键改装推荐</h2>
        <p>选择枪械和偏好风格，系统从改枪码库 + 配件库中智能推荐最优方案</p>
      </div>

      <div class="smart-steps">
        <div class="step-panel">
          <h3><span class="step-num">1</span> 选择枪械</h3>
          <select class="custom-select" @change="selectWeapon">
            <option value="">-- 选择一把枪械 --</option>
            <option v-for="w in weapons" :key="w.id" :value="w.id">
              {{ w.name }} ({{ w.type }}) - {{ w.tier }}
            </option>
          </select>

          <div v-if="selectedWeapon" style="margin-top:16px">
            <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:8px">
              <span style="font-weight:600">{{ selectedWeapon.name }}</span>
              <span class="tier-badge" :class="'tier-'+selectedWeapon.tier">{{ selectedWeapon.tier }}</span>
            </div>
            <div class="card-meta">
              <span class="type-icon" :class="'type-'+selectedWeapon.type">{{ selectedWeapon.type }}</span>
              <span>{{ selectedWeapon.caliber }}</span>
              <span>伤害 {{ selectedWeapon.base_damage }}</span>
              <span>射速 {{ selectedWeapon.max_rpm }}</span>
            </div>
          </div>

          <h3 style="margin-top:20px"><span class="step-num">2</span> 选择偏好风格</h3>
          <div class="style-buttons">
            <button v-for="s in styles" :key="s.key"
                    class="style-btn" :class="{active: style===s.key}"
                    @click="setStyle(s.key)">
              {{ s.label }}
            </button>
          </div>

          <h3 style="margin-top:8px"><span class="step-num">3</span> 预算设置</h3>
          <div class="budget-input">
            <label>预算上限</label>
            <input type="number" v-model.number="budget" step="50000" min="0" max="2000000">
            <span class="price-tag">{{ formatPrice(budget) }}</span>
          </div>
        </div>

        <div class="step-panel">
          <h3><span class="step-num">4</span> 调整目标属性</h3>
          <p style="font-size:12px;color:var(--text-muted);margin-bottom:12px">
            拖动滑块调整你期望的属性数值（左侧为裸枪基础值）
          </p>
          <div v-for="attr in attrLabels" :key="attr.key" class="slider-group">
            <div class="slider-label">
              <span>{{ attr.label }}</span>
              <span>{{ attrs[attr.key] }}</span>
            </div>
            <input type="range" min="0" max="100" v-model.number="attrs[attr.key]">
          </div>

          <button class="action-btn" @click="doRecommend" :disabled="!selectedWeapon || loading">
            {{ loading ? '智能分析中...' : '开始智能推荐' }}
          </button>
          <p v-if="!selectedWeapon" style="text-align:center;font-size:12px;color:var(--text-muted);margin-top:6px">
            请先在左侧选择一把枪械
          </p>
        </div>
      </div>

      <div v-if="error" class="results-panel" style="border-color:#ef4444">
        <p style="color:#ef4444;text-align:center">{{ error }}</p>
      </div>

      <div v-if="searched && !error" class="results-panel">
        <h3 style="color:var(--accent);margin-bottom:8px;font-size:18px">
          推荐结果
          <span style="font-size:13px;color:var(--text-muted);font-weight:400">
            (共{{ results.length }}个方案)
          </span>
        </h3>
        <p style="font-size:12px;color:var(--text-muted);margin-bottom:16px">
          预算 {{ formatPrice(budget) }} → 丐版(预算30%) / 半改(预算70%) / 满改(预算100%) 三档自动生成
        </p>

        <div v-if="results.length===0" class="empty-state">
          暂无符合条件的改装方案，请尝试提高预算或调整属性偏好
        </div>

        <div v-for="(r, i) in results" :key="i" class="result-item"
             :style="tierBorder(r.mod_code.grade)">
          <div class="result-header">
            <div>
              <div style="display:flex;align-items:center;gap:8px;flex-wrap:wrap">
                <span class="grade-badge" :class="'grade-'+r.mod_code.grade" style="font-size:13px;padding:3px 12px">
                  {{ r.mod_code.grade }}
                </span>
                <span style="font-weight:700;font-size:16px">{{ r.mod_code.name }}</span>
                <span v-if="isDynamic(r)" class="tag" style="background:rgba(0,194,255,0.15);color:var(--accent);border-color:var(--accent)">
                  AI配件组合
                </span>
              </div>
              <div class="card-meta" style="margin:6px 0">
                <span class="price-tag" style="font-size:15px">{{ formatPrice(r.mod_code.total_price) }}</span>
                <span>{{ r.mod_code.weapon_name }}</span>
                <span style="color:var(--text-muted)">配件{{ countParts(r) }}件</span>
              </div>
            </div>
            <div class="match-score" :style="{color: scoreColor(r.match_score)}">{{ r.match_score }}%</div>
          </div>

          <div class="code-block">
            <span class="code-text">{{ r.mod_code.code }}</span>
            <button v-if="!isDynamic(r)" class="copy-btn" :class="{copied: copiedId===i}"
                    @click="copyCode(r.mod_code.code, i)">
              {{ copiedId===i ? '已复制' : '一键复制' }}
            </button>
          </div>

          <div style="display:flex;gap:20px;margin-top:12px;flex-wrap:wrap">
            <div class="radar-wrap" style="flex:0 0 auto">
              <radar-chart :datasets="resultRadar(r)" :size="200"></radar-chart>
            </div>
            <div style="flex:1;min-width:240px">
              <div style="font-size:12px;color:var(--text-muted);margin-bottom:8px">
                <strong>配件清单：</strong>{{ r.mod_code.parts }}
              </div>

              <div style="margin:8px 0">
                <div v-for="attr in attrLabels" :key="'d'+attr.key"
                     style="display:flex;align-items:center;gap:6px;font-size:12px;margin-bottom:3px">
                  <span style="width:70px;color:var(--text-secondary)">{{ attr.label }}</span>
                  <div class="stat-bar-wrap" style="flex:1;height:5px">
                    <div class="stat-bar" :style="{width: r.mod_code[attr.key]+'%'}"></div>
                  </div>
                  <span style="width:28px;text-align:right;font-family:'Rajdhani';font-weight:600">{{ Math.round(r.mod_code[attr.key]) }}</span>
                  <span :class="diffClass(attrDiff(r, attr.key))"
                        style="width:36px;text-align:right;font-size:11px;font-weight:600">
                    {{ diffStr(attrDiff(r, attr.key)) }}
                  </span>
                </div>
              </div>

              <div class="tags-wrap" v-if="r.mod_code.tags">
                <span class="tag" v-for="t in r.mod_code.tags.split(',')" :key="t">{{ t.trim() }}</span>
              </div>
              <div style="font-size:13px;color:var(--text-secondary);margin-top:6px;line-height:1.5">
                {{ r.reason }}
              </div>
              <div class="card-desc" style="white-space:pre-wrap">{{ r.mod_code.description }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `
};
