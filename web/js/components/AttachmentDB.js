const AttachmentDB = {
  data() {
    return {
      attachments: [],
      loading: true,
      filterSlot: ''
    };
  },
  async created() {
    await this.load();
  },
  computed: {
    slots() {
      return ['', '枪口', '枪管', '前握把', '后握把', '枪托', '弹匣', '瞄具', '导轨', '导气'];
    },
    grouped() {
      if (this.filterSlot) return { [this.filterSlot]: this.attachments };
      const g = {};
      this.attachments.forEach(a => {
        if (!g[a.slot]) g[a.slot] = [];
        g[a.slot].push(a);
      });
      return g;
    }
  },
  methods: {
    async load() {
      this.loading = true;
      try {
        const params = {};
        if (this.filterSlot) params.slot = this.filterSlot;
        this.attachments = await API.getAttachments(params);
      } catch(e) { console.error(e); }
      this.loading = false;
    },
    setSlot(s) { this.filterSlot = s; this.load(); },
    fmtVal(v) {
      if (v > 0) return '+' + v;
      if (v < 0) return '' + v;
      return '-';
    },
    valClass(v) {
      if (v > 0) return 'positive';
      if (v < 0) return 'negative';
      return '';
    },
    formatPrice(p) {
      if (p >= 10000) return (p / 10000).toFixed(1) + '万';
      return p.toLocaleString();
    }
  },
  template: `
    <div>
      <div class="section-header">
        <h2><span class="icon">&#9881;</span> 配件数据库</h2>
        <p>全部配件属性加成一览，绿色为正向加成，红色为负向影响</p>
      </div>

      <div class="filters-bar">
        <button v-for="s in slots" :key="s"
                class="filter-btn" :class="{active: filterSlot===s}"
                @click="setSlot(s)">
          {{ s || '全部' }}
        </button>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-else>
        <div v-for="(items, slot) in grouped" :key="slot" style="margin-bottom:24px">
          <h3 style="font-size:16px;margin-bottom:12px;display:flex;align-items:center;gap:8px">
            <span :class="'slot-'+slot" style="font-weight:700">{{ slot }}</span>
            <span style="font-size:12px;color:var(--text-muted)">{{ items.length }}件</span>
          </h3>
          <div style="overflow-x:auto">
            <table class="data-table">
              <thead>
                <tr>
                  <th>名称</th>
                  <th>价格</th>
                  <th>射程</th>
                  <th>垂直后坐</th>
                  <th>水平后坐</th>
                  <th>操控</th>
                  <th>据枪</th>
                  <th>腰射</th>
                  <th>初速</th>
                  <th>消音</th>
                  <th>精校</th>
                  <th>适用</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="a in items" :key="a.id">
                  <td style="font-weight:600">{{ a.name }}</td>
                  <td class="price-tag">{{ formatPrice(a.price) }}</td>
                  <td :class="valClass(a.effective_range)">{{ fmtVal(a.effective_range) }}</td>
                  <td :class="valClass(a.vertical_recoil)">{{ fmtVal(a.vertical_recoil) }}</td>
                  <td :class="valClass(a.horiz_recoil)">{{ fmtVal(a.horiz_recoil) }}</td>
                  <td :class="valClass(a.handling_speed)">{{ fmtVal(a.handling_speed) }}</td>
                  <td :class="valClass(a.ads_stability)">{{ fmtVal(a.ads_stability) }}</td>
                  <td :class="valClass(a.hip_fire_acc)">{{ fmtVal(a.hip_fire_acc) }}</td>
                  <td :class="valClass(a.muzzle_velocity)">{{ fmtVal(a.muzzle_velocity) }}</td>
                  <td :class="valClass(a.sound_range)">{{ fmtVal(a.sound_range) }}</td>
                  <td style="font-size:11px;color:var(--text-muted)">
                    <template v-if="a.tune_attr_a">{{ a.tune_attr_a }} / {{ a.tune_attr_b }}</template>
                    <template v-else>-</template>
                  </td>
                  <td style="font-size:11px;color:var(--text-muted)">{{ a.compat_weapons }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  `
};
