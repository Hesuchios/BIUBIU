const WeaponDB = {
  data() {
    return {
      weapons: [],
      loading: true,
      filterType: '',
      filterTier: '',
      sortBy: 'rank',
      compareIds: [],
      compareData: null,
      showCompare: false,
      expanded: null
    };
  },
  async created() {
    await this.load();
  },
  computed: {
    types() { return ['', 'AR', 'SMG', 'SR', 'DMR', 'LMG', 'SG', 'Pistol']; },
    tiers() { return ['', 'T0', 'T1', 'T2', 'T3']; },
    typeLabels() { return {'': '全部', AR:'步枪', SMG:'冲锋枪', SR:'狙击枪', DMR:'射手步枪', LMG:'轻机枪', SG:'霰弹枪', Pistol:'手枪'}; },
  },
  methods: {
    async load() {
      this.loading = true;
      try {
        const params = {};
        if (this.filterType) params.type = this.filterType;
        if (this.filterTier) params.tier = this.filterTier;
        params.sort = this.sortBy;
        this.weapons = await API.getWeapons(params);
      } catch(e) { console.error(e); }
      this.loading = false;
    },
    formatPrice(p) {
      if (p >= 10000) return (p / 10000).toFixed(1) + '万';
      return p.toLocaleString();
    },
    setType(t) { this.filterType = t; this.load(); },
    setTier(t) { this.filterTier = t; this.load(); },
    setSort(s) { this.sortBy = s; this.load(); },
    toggleExpand(id) { this.expanded = this.expanded === id ? null : id; },
    toggleCompare(id) {
      const idx = this.compareIds.indexOf(id);
      if (idx >= 0) {
        this.compareIds.splice(idx, 1);
      } else if (this.compareIds.length < 3) {
        this.compareIds.push(id);
      }
    },
    async doCompare() {
      if (this.compareIds.length < 2) return;
      try {
        this.compareData = await API.compareWeapons(this.compareIds);
        this.showCompare = true;
      } catch(e) { console.error(e); }
    },
    radarValues(w) {
      return [{
        name: w.name,
        values: [w.effective_range, w.vertical_recoil, w.horiz_recoil,
                 w.handling_speed, w.ads_stability, w.hip_fire_acc,
                 w.muzzle_velocity, w.sound_range]
      }];
    },
    compareRadar() {
      if (!this.compareData) return [];
      return this.compareData.map(w => ({
        name: w.name,
        values: [w.effective_range, w.vertical_recoil, w.horiz_recoil,
                 w.handling_speed, w.ads_stability, w.hip_fire_acc,
                 w.muzzle_velocity, w.sound_range]
      }));
    }
  },
  template: `
    <div>
      <div class="section-header">
        <h2><span class="icon">&#9776;</span> 枪械数据库</h2>
        <p>全部 {{ weapons.length }} 把枪械完整属性数据</p>
      </div>

      <div class="filters-bar">
        <button v-for="t in types" :key="'t'+t"
                class="filter-btn" :class="{active: filterType===t}"
                @click="setType(t)">
          {{ typeLabels[t] || t }}
        </button>
        <span style="width:1px;background:var(--border);margin:0 4px"></span>
        <button v-for="t in tiers" :key="'ti'+t"
                class="filter-btn" :class="{active: filterTier===t}"
                @click="setTier(t)">
          {{ t || '全梯队' }}
        </button>
        <span style="width:1px;background:var(--border);margin:0 4px"></span>
        <select class="custom-select" style="width:auto;padding:4px 12px;font-size:12px" @change="setSort($event.target.value)">
          <option value="rank">按排名</option>
          <option value="damage">按伤害</option>
          <option value="rpm">按射速</option>
          <option value="range">按射程</option>
          <option value="price">按价格</option>
          <option value="name">按名称</option>
        </select>
      </div>

      <div v-if="compareIds.length >= 2" style="margin-bottom:16px">
        <button class="action-btn" style="width:auto;padding:8px 20px" @click="doCompare">
          对比选中 ({{ compareIds.length }}把)
        </button>
      </div>

      <div v-if="showCompare && compareData" class="compare-section">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
          <h3 style="margin:0;color:var(--accent)">枪械对比</h3>
          <button class="filter-btn" @click="showCompare=false;compareIds=[];compareData=null">关闭</button>
        </div>
        <div class="radar-wrap">
          <radar-chart :datasets="compareRadar()" :size="300"></radar-chart>
        </div>
        <div style="display:flex;gap:16px;justify-content:center;margin-bottom:12px">
          <span v-for="(w,i) in compareData" :key="w.id"
            :style="{color: ['#00c2ff','#ff6b35','#32cd32'][i]}">
            &#9632; {{ w.name }}
          </span>
        </div>
        <table class="data-table">
          <thead>
            <tr>
              <th>属性</th>
              <th v-for="w in compareData" :key="w.id">{{ w.name }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="attr in [
              {label:'基础伤害',key:'base_damage'},{label:'护甲伤害',key:'armor_damage'},
              {label:'射速RPM',key:'max_rpm'},{label:'射程',key:'effective_range'},
              {label:'垂直后坐',key:'vertical_recoil'},{label:'水平后坐',key:'horiz_recoil'},
              {label:'操控速度',key:'handling_speed'},{label:'据枪稳定',key:'ads_stability'},
              {label:'腰射精度',key:'hip_fire_acc'},{label:'枪口初速',key:'muzzle_velocity'},
              {label:'消音能力',key:'sound_range'}
            ]" :key="attr.key">
              <td>{{ attr.label }}</td>
              <td v-for="w in compareData" :key="w.id">{{ w[attr.key] }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="weapons.length===0" class="empty-state">暂无数据</div>
      <div v-else class="cards-grid">
        <div v-for="w in weapons" :key="w.id" class="weapon-card">
          <div class="card-top" @click="toggleExpand(w.id)">
            <div>
              <div class="card-name">{{ w.name }}</div>
              <div class="card-meta">
                <span class="type-icon" :class="'type-'+w.type">{{ w.type }}</span>
                <span>{{ w.caliber }}</span>
                <span class="price-tag">{{ formatPrice(w.base_price) }}</span>
              </div>
            </div>
            <div style="display:flex;flex-direction:column;align-items:flex-end;gap:6px">
              <span class="tier-badge" :class="'tier-'+w.tier">{{ w.tier }}</span>
              <label @click.stop style="display:flex;align-items:center;gap:4px;font-size:11px;color:var(--text-muted);cursor:pointer">
                <input type="checkbox" :checked="compareIds.includes(w.id)"
                       @change="toggleCompare(w.id)" @click.stop>
                对比
              </label>
            </div>
          </div>

          <div class="card-stats" @click="toggleExpand(w.id)">
            <div class="stat-row"><span class="stat-label">伤害</span><span class="stat-value">{{ w.base_damage }}</span></div>
            <div class="stat-row"><span class="stat-label">甲伤</span><span class="stat-value">{{ w.armor_damage }}</span></div>
            <div class="stat-row"><span class="stat-label">射速</span><span class="stat-value">{{ w.max_rpm }}</span></div>
            <div class="stat-row"><span class="stat-label">模式</span><span class="stat-value" style="font-size:12px">{{ w.fire_mode }}</span></div>
          </div>

          <div class="card-expand" v-if="expanded === w.id">
            <div class="radar-wrap">
              <radar-chart :datasets="radarValues(w)" :size="220"></radar-chart>
            </div>
            <div class="stat-row-full" v-for="stat in [
              {label:'优势射程', key:'effective_range'},
              {label:'垂直后坐', key:'vertical_recoil'},
              {label:'水平后坐', key:'horiz_recoil'},
              {label:'操控速度', key:'handling_speed'},
              {label:'据枪稳定', key:'ads_stability'},
              {label:'腰射精度', key:'hip_fire_acc'},
              {label:'枪口初速', key:'muzzle_velocity'},
              {label:'消音能力', key:'sound_range'}
            ]" :key="stat.key">
              <span class="stat-label">{{ stat.label }}</span>
              <div class="stat-bar-wrap">
                <div class="stat-bar" :style="{width: w[stat.key]+'%'}"></div>
              </div>
              <span class="stat-value">{{ w[stat.key] }}</span>
            </div>
            <div class="card-desc" style="margin-top:12px">{{ w.description }}</div>
          </div>
        </div>
      </div>
    </div>
  `
};
