const KnowledgeBase = {
  data() {
    return {
      kb: null,
      loading: true,
      activeTab: 'principles'
    };
  },
  async created() {
    try {
      this.kb = await API.getKnowledge();
    } catch(e) { console.error(e); }
    this.loading = false;
  },
  computed: {
    tabs() {
      return [
        { key: 'principles', label: '核心原则' },
        { key: 'attributes', label: '属性详解' },
        { key: 'slots', label: '配件选择' },
        { key: 'tuning', label: '精校系统' },
        { key: 'distance', label: '距离策略' },
        { key: 'weapon_types', label: '武器类型' },
        { key: 'budget', label: '预算分档' },
        { key: 'mistakes', label: '避坑指南' },
        { key: 's9', label: 'S9新内容' }
      ];
    }
  },
  methods: {
    setTab(t) { this.activeTab = t; }
  },
  template: `
    <div>
      <div class="section-header">
        <h2><span class="icon">&#128218;</span> 改枪知识库</h2>
        <p>基于社区攻略、实战数据和机制研究提炼的改枪核心知识 ({{ kb ? kb.version + '赛季' : '加载中...' }})</p>
      </div>

      <div v-if="loading" class="loading">加载知识库中...</div>

      <div v-if="kb && !loading">
        <div class="filters-bar">
          <button v-for="t in tabs" :key="t.key"
                  class="filter-btn" :class="{active: activeTab===t.key}"
                  @click="setTab(t.key)">
            {{ t.label }}
          </button>
        </div>

        <!-- 核心原则 -->
        <div v-if="activeTab==='principles'" class="kb-section">
          <div class="kb-card">
            <h3>{{ kb.core_principles.title }}</h3>
            <ul class="kb-list">
              <li v-for="(rule, i) in kb.core_principles.rules" :key="i">{{ rule }}</li>
            </ul>
          </div>
        </div>

        <!-- 属性详解 -->
        <div v-if="activeTab==='attributes'" class="kb-section">
          <div v-for="attr in kb.attribute_guide.attributes" :key="attr.key" class="kb-card">
            <div class="kb-card-header">
              <h3>{{ attr.name }}</h3>
              <span class="tag" style="background:rgba(0,194,255,0.1);color:var(--accent);border-color:var(--accent)">
                {{ attr.importance }}
              </span>
            </div>
            <p class="kb-desc">{{ attr.description }}</p>
            <div class="kb-attr-grid">
              <div>
                <span class="kb-label positive">提升配件</span>
                <span v-for="p in attr.improve_by" :key="p" class="tag">{{ p }}</span>
              </div>
              <div>
                <span class="kb-label negative">降低配件</span>
                <span v-for="p in attr.reduce_by" :key="p" class="tag" style="border-color:rgba(239,68,68,0.3);color:#ef4444">{{ p }}</span>
              </div>
            </div>
            <div class="kb-tip">{{ attr.tip }}</div>
          </div>
        </div>

        <!-- 配件选择 -->
        <div v-if="activeTab==='slots'" class="kb-section">
          <div v-for="s in kb.slot_guide.slots" :key="s.slot" class="kb-card">
            <div class="kb-card-header">
              <h3 :class="'slot-'+s.slot">{{ s.slot }}</h3>
              <span class="tag">{{ s.role }}</span>
              <span class="tag" :style="s.priority==='最高'||s.priority==='必装'?'background:rgba(255,215,0,0.15);color:#ffd700;border-color:rgba(255,215,0,0.3)':''">
                优先级: {{ s.priority }}
              </span>
            </div>
            <ul class="kb-list">
              <li v-for="(rule, i) in s.rules" :key="i">{{ rule }}</li>
            </ul>
          </div>
        </div>

        <!-- 精校系统 -->
        <div v-if="activeTab==='tuning'" class="kb-section">
          <div class="kb-card" style="border-color:rgba(0,194,255,0.3)">
            <h3 style="color:var(--accent)">{{ kb.tuning_guide.universal_formula.name }}</h3>
            <div class="kb-formula">{{ kb.tuning_guide.universal_formula.formula }}</div>
            <p class="kb-desc">{{ kb.tuning_guide.universal_formula.explanation }}</p>
            <div style="margin-top:12px">
              <div v-for="(desc, name) in kb.tuning_guide.universal_formula.variants" :key="name"
                   style="margin-bottom:8px">
                <span class="tag" style="background:rgba(168,85,247,0.15);color:#a855f7;border-color:rgba(168,85,247,0.3)">{{ name }}</span>
                <span style="font-size:13px;color:var(--text-secondary);margin-left:8px">{{ desc }}</span>
              </div>
            </div>
          </div>
          <div class="kb-card">
            <h3>各部位精校参数</h3>
            <p class="kb-desc">{{ kb.tuning_guide.overview }}</p>
            <div v-for="pt in kb.tuning_guide.parts_tuning" :key="pt.part" style="margin-top:12px;padding:12px;background:rgba(0,0,0,0.2);border-radius:8px">
              <div style="font-weight:600;margin-bottom:6px">{{ pt.part }}</div>
              <div style="font-size:12px;color:var(--text-secondary)">
                <div>参数: {{ pt.params.join(' / ') }}</div>
                <div v-if="pt.formula" style="color:var(--accent)">通用公式方向: {{ pt.formula }}</div>
                <div v-if="pt['远程']">远程: {{ pt['远程'] }}</div>
                <div v-if="pt['近战']">近战: {{ pt['近战'] }}</div>
                <div v-if="pt.tip">{{ pt.tip }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 距离策略 -->
        <div v-if="activeTab==='distance'" class="kb-section">
          <div v-for="s in kb.distance_strategy.strategies" :key="s.range" class="kb-card">
            <h3>{{ s.range }}</h3>
            <div class="kb-tip" style="margin-bottom:10px">优先级: {{ s.priority }}</div>
            <div class="kb-config-grid">
              <div v-for="(val, slot) in s.config" :key="slot" class="kb-config-item">
                <span class="kb-label" :class="'slot-'+slot">{{ slot }}</span>
                <span>{{ val }}</span>
              </div>
            </div>
            <div class="kb-tip" style="margin-top:8px">{{ s.tip }}</div>
          </div>
        </div>

        <!-- 武器类型 -->
        <div v-if="activeTab==='weapon_types'" class="kb-section">
          <div v-for="t in kb.weapon_type_guide.types" :key="t.type" class="kb-card">
            <div class="kb-card-header">
              <h3><span class="type-icon" :class="'type-'+t.type">{{ t.type }}</span> {{ t.name }}</h3>
            </div>
            <p class="kb-desc">{{ t.guide }}</p>
            <div class="kb-tip">{{ t.key_point }}</div>
            <div v-if="t.examples.length" style="margin-top:6px">
              <span v-for="ex in t.examples" :key="ex" class="tag" style="margin-right:4px;margin-bottom:4px">{{ ex }}</span>
            </div>
          </div>
        </div>

        <!-- 预算分档 -->
        <div v-if="activeTab==='budget'" class="kb-section">
          <div v-for="t in kb.budget_tiers.tiers" :key="t.grade" class="kb-card">
            <div class="kb-card-header">
              <h3><span class="grade-badge" :class="'grade-'+t.grade">{{ t.grade }}</span> {{ t.budget_range }}</h3>
            </div>
            <p class="kb-desc">{{ t.description }}</p>
            <div class="kb-tip">{{ t.strategy }}</div>
          </div>
        </div>

        <!-- 避坑指南 -->
        <div v-if="activeTab==='mistakes'" class="kb-section">
          <div class="kb-card" style="border-color:rgba(239,68,68,0.3)">
            <h3 style="color:#ef4444">{{ kb.common_mistakes.title }}</h3>
            <ul class="kb-list kb-mistakes">
              <li v-for="(m, i) in kb.common_mistakes.mistakes" :key="i">{{ m }}</li>
            </ul>
          </div>
        </div>

        <!-- S9新内容 -->
        <div v-if="activeTab==='s9'" class="kb-section">
          <div class="kb-card" style="border-color:rgba(255,215,0,0.3)">
            <h3 style="color:var(--accent-gold)">{{ kb.s9_highlights.title }}</h3>
            <ul class="kb-list">
              <li v-for="(item, i) in kb.s9_highlights.items" :key="i">{{ item }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  `
};
