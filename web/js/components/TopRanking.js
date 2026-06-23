const TopRanking = {
  data() {
    return {
      weapons: [],
      loading: true,
      expanded: null,
      grade: 'all'
    };
  },
  computed: {
    grades() {
      return [
        { key: 'all', label: '全部推荐' },
        { key: 'budget', label: '丐版 (总价最低)' },
        { key: 'half', label: '半改 (性价比)' },
        { key: 'full', label: '满改 (顶配)' }
      ];
    }
  },
  async created() {
    await this.loadData();
  },
  methods: {
    async loadData() {
      this.loading = true;
      try {
        this.weapons = await API.get('/api/weapons/top10?grade=' + this.grade);
      } catch (e) {
        console.error(e);
      }
      this.loading = false;
    },
    async setGrade(g) {
      this.grade = g;
      this.expanded = null;
      await this.loadData();
    },
    toggle(id) {
      this.expanded = this.expanded === id ? null : id;
    },
    rankClass(i) {
      if (i === 0) return 'gold';
      if (i === 1) return 'silver';
      if (i === 2) return 'bronze';
      return '';
    },
    formatPrice(p) {
      if (p >= 10000) return (p / 10000).toFixed(1) + '万';
      return p.toLocaleString();
    },
    gradeClass(tag) {
      if (tag === '丐版') return 'grade-budget';
      if (tag === '半改') return 'grade-half';
      if (tag === '满改') return 'grade-full';
      return '';
    },
    radarValues(w) {
      return [{
        name: w.name,
        values: [w.effective_range, w.vertical_recoil, w.horiz_recoil,
                 w.handling_speed, w.ads_stability, w.hip_fire_acc,
                 w.muzzle_velocity, w.sound_range]
      }];
    }
  },
  template: `
    <div>
      <div class="section-header">
        <h2><span class="icon">&#9733;</span> TOP 10 枪械推荐榜</h2>
        <p>基于 S9 赛季综合强度、性价比、实战表现评选（价格为交易行实际成交价）</p>
      </div>
      <div class="filters-bar">
        <button v-for="g in grades" :key="g.key"
                class="filter-btn" :class="{active: grade===g.key}"
                @click="setGrade(g.key)">
          {{ g.label }}
        </button>
      </div>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else class="cards-grid">
        <div v-for="(w, i) in weapons" :key="w.id"
             class="weapon-card" @click="toggle(w.id)">
          <div class="card-top">
            <span class="card-rank" :class="rankClass(i)">#{{ i + 1 }}</span>
            <span class="tier-badge" :class="'tier-'+w.tier">{{ w.tier }}</span>
          </div>
          <div class="card-name">{{ w.name }}</div>
          <div class="card-meta">
            <span class="type-icon" :class="'type-'+w.type">{{ w.type }}</span>
            <span>{{ w.caliber }}</span>
            <span>{{ w.fire_mode }}</span>
            <span class="price-tag">裸枪 {{ formatPrice(w.base_price) }}</span>
          </div>
          <div class="card-cost-row" v-if="w.total_cost">
            <span class="grade-badge" :class="gradeClass(w.grade_tag)">{{ w.grade_tag }}</span>
            <span class="total-cost">总花费 <strong>{{ formatPrice(w.total_cost) }}</strong></span>
            <span class="build-tip" v-if="w.build_tip">{{ w.build_tip }}</span>
          </div>

          <div class="card-stats">
            <div class="stat-row">
              <span class="stat-label">伤害</span>
              <span class="stat-value">{{ w.base_damage }}</span>
            </div>
            <div class="stat-row">
              <span class="stat-label">甲伤</span>
              <span class="stat-value">{{ w.armor_damage }}</span>
            </div>
            <div class="stat-row">
              <span class="stat-label">射速</span>
              <span class="stat-value">{{ w.max_rpm }}RPM</span>
            </div>
            <div class="stat-row">
              <span class="stat-label">射程</span>
              <span class="stat-value">{{ w.effective_range }}</span>
            </div>
          </div>

          <div class="card-desc">{{ w.description }}</div>

          <div class="card-expand" v-if="expanded === w.id" @click.stop>
            <div class="radar-wrap">
              <radar-chart :datasets="radarValues(w)" :size="240"></radar-chart>
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
          </div>
        </div>
      </div>
    </div>
  `
};
