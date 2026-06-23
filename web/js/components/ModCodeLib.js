const ModCodeLib = {
  data() {
    return {
      codes: [],
      weapons: [],
      loading: true,
      filterWeapon: '',
      filterGrade: '',
      filterTag: '',
      copiedId: null
    };
  },
  async created() {
    try {
      const [codes, weapons] = await Promise.all([
        API.getModCodes(),
        API.getWeapons()
      ]);
      this.codes = codes;
      this.weapons = weapons;
    } catch(e) { console.error(e); }
    this.loading = false;
  },
  computed: {
    weaponNames() {
      const names = new Set(this.codes.map(c => c.weapon_name));
      return ['', ...names];
    },
    grades() { return ['', '丐版', '半改', '满改']; },
    tags() { return ['', '近战', '中距', '远程', '性价比', '稳定', '腰射', '消音', '架点']; },
    filtered() {
      let list = this.codes;
      if (this.filterWeapon) list = list.filter(c => c.weapon_name === this.filterWeapon);
      if (this.filterGrade) list = list.filter(c => c.grade === this.filterGrade);
      if (this.filterTag) list = list.filter(c => c.tags && c.tags.includes(this.filterTag));
      return list;
    },
    grouped() {
      const g = {};
      this.filtered.forEach(c => {
        if (!g[c.weapon_name]) g[c.weapon_name] = [];
        g[c.weapon_name].push(c);
      });
      return g;
    }
  },
  methods: {
    formatPrice(p) {
      if (p >= 10000) return (p / 10000).toFixed(1) + '万';
      return p.toLocaleString();
    },
    async copyCode(code, id) {
      try {
        await navigator.clipboard.writeText(code);
        this.copiedId = id;
        setTimeout(() => { this.copiedId = null; }, 2000);
      } catch(e) {
        const ta = document.createElement('textarea');
        ta.value = code;
        document.body.appendChild(ta);
        ta.select();
        document.execCommand('copy');
        document.body.removeChild(ta);
        this.copiedId = id;
        setTimeout(() => { this.copiedId = null; }, 2000);
      }
    },
    getWeaponForCode(mc) {
      return this.weapons.find(w => w.id === mc.weapon_id);
    },
    radarData(mc) {
      const ds = [{
        name: '改装后',
        values: [mc.effective_range, mc.vertical_recoil, mc.horiz_recoil,
                 mc.handling_speed, mc.ads_stability, mc.hip_fire_acc,
                 mc.muzzle_velocity, mc.sound_range]
      }];
      const base = this.getWeaponForCode(mc);
      if (base) {
        ds.unshift({
          name: '裸枪',
          values: [base.effective_range, base.vertical_recoil, base.horiz_recoil,
                   base.handling_speed, base.ads_stability, base.hip_fire_acc,
                   base.muzzle_velocity, base.sound_range]
        });
      }
      return ds;
    }
  },
  template: `
    <div>
      <div class="section-header">
        <h2><span class="icon">&#128273;</span> 改枪码库</h2>
        <p>收录各武器改枪码方案，点击即可一键复制</p>
      </div>

      <div class="filters-bar">
        <select class="custom-select" style="width:auto;padding:4px 12px;font-size:12px"
                @change="filterWeapon=$event.target.value">
          <option value="">全部武器</option>
          <option v-for="n in weaponNames.slice(1)" :key="n" :value="n">{{ n }}</option>
        </select>
        <button v-for="g in grades" :key="'g'+g"
                class="filter-btn" :class="{active: filterGrade===g}"
                @click="filterGrade=g">
          {{ g || '全部档次' }}
        </button>
        <span style="width:1px;background:var(--border);margin:0 4px"></span>
        <button v-for="t in tags" :key="'tg'+t"
                class="filter-btn" :class="{active: filterTag===t}"
                @click="filterTag=t">
          {{ t || '全部标签' }}
        </button>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="filtered.length===0" class="empty-state">暂无匹配的改枪码</div>
      <div v-else>
        <div v-for="(items, weaponName) in grouped" :key="weaponName" style="margin-bottom:28px">
          <h3 style="font-size:17px;font-weight:700;margin-bottom:14px;padding-bottom:8px;border-bottom:1px solid var(--border)">
            {{ weaponName }}
            <span style="font-size:12px;color:var(--text-muted);font-weight:400;margin-left:8px">{{ items.length }}个方案</span>
          </h3>
          <div class="cards-grid">
            <div v-for="mc in items" :key="mc.id" class="weapon-card" style="cursor:default">
              <div class="card-top">
                <div>
                  <div class="card-name" style="font-size:15px">{{ mc.name }}</div>
                  <div class="card-meta" style="margin-bottom:8px">
                    <span class="grade-badge" :class="'grade-'+mc.grade">{{ mc.grade }}</span>
                    <span class="price-tag">{{ formatPrice(mc.total_price) }}</span>
                  </div>
                </div>
              </div>

              <div class="code-block">
                <span class="code-text">{{ mc.code }}</span>
                <button class="copy-btn" :class="{copied: copiedId===mc.id}"
                        @click="copyCode(mc.code, mc.id)">
                  {{ copiedId===mc.id ? '已复制' : '复制' }}
                </button>
              </div>

              <div class="tags-wrap" v-if="mc.tags">
                <span class="tag" v-for="t in mc.tags.split(',')" :key="t">{{ t }}</span>
              </div>

              <div style="font-size:12px;color:var(--text-muted);margin-top:8px">
                <strong>配件：</strong>{{ mc.parts }}
              </div>

              <div class="radar-wrap" style="margin-top:10px">
                <radar-chart :datasets="radarData(mc)" :size="200"></radar-chart>
              </div>

              <div class="card-desc">{{ mc.description }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `
};
