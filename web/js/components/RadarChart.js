const RadarChart = {
  props: {
    datasets: { type: Array, required: true },
    size: { type: Number, default: 220 },
    labels: {
      type: Array,
      default: () => ['射程', '垂直稳定', '水平稳定', '操控', '据枪', '腰射', '初速', '消音']
    }
  },
  computed: {
    center() { return this.size / 2; },
    radius() { return this.size / 2 - 30; },
    axes() {
      const n = this.labels.length;
      return this.labels.map((label, i) => {
        const angle = (Math.PI * 2 * i) / n - Math.PI / 2;
        return {
          x: this.center + this.radius * Math.cos(angle),
          y: this.center + this.radius * Math.sin(angle),
          lx: this.center + (this.radius + 18) * Math.cos(angle),
          ly: this.center + (this.radius + 18) * Math.sin(angle),
          label, angle
        };
      });
    },
    gridLevels() {
      return [0.2, 0.4, 0.6, 0.8, 1.0];
    },
    gridPolygons() {
      const n = this.labels.length;
      return this.gridLevels.map(level => {
        const pts = [];
        for (let i = 0; i < n; i++) {
          const angle = (Math.PI * 2 * i) / n - Math.PI / 2;
          pts.push(
            (this.center + this.radius * level * Math.cos(angle)) + ',' +
            (this.center + this.radius * level * Math.sin(angle))
          );
        }
        return pts.join(' ');
      });
    },
    polygons() {
      const colors = ['rgba(0,194,255,0.35)', 'rgba(255,107,53,0.30)', 'rgba(50,205,50,0.30)'];
      const strokes = ['#00c2ff', '#ff6b35', '#32cd32'];
      const n = this.labels.length;
      return this.datasets.map((ds, di) => {
        const points = ds.values.map((v, i) => {
          const pct = Math.min(v, 100) / 100;
          const angle = (Math.PI * 2 * i) / n - Math.PI / 2;
          return {
            x: this.center + this.radius * pct * Math.cos(angle),
            y: this.center + this.radius * pct * Math.sin(angle)
          };
        });
        return {
          pointsStr: points.map(p => p.x + ',' + p.y).join(' '),
          fill: colors[di % colors.length],
          stroke: strokes[di % strokes.length],
          name: ds.name,
          points
        };
      });
    }
  },
  template: `
    <svg :width="size" :height="size" class="radar-chart">
      <polygon v-for="(gp, gi) in gridPolygons" :key="'g'+gi"
        :points="gp" fill="none" stroke="rgba(255,255,255,0.08)" stroke-width="1"/>
      <line v-for="a in axes" :key="'l'+a.label"
        :x1="center" :y1="center" :x2="a.x" :y2="a.y"
        stroke="rgba(255,255,255,0.12)" stroke-width="1"/>
      <g v-for="(poly, pi) in polygons" :key="'p'+pi">
        <polygon :points="poly.pointsStr"
          :fill="poly.fill" :stroke="poly.stroke" stroke-width="2"/>
        <circle v-for="(pt, pti) in poly.points" :key="pti"
          :cx="pt.x" :cy="pt.y" r="3" :fill="poly.stroke"/>
      </g>
      <text v-for="a in axes" :key="'t'+a.label"
        :x="a.lx" :y="a.ly"
        text-anchor="middle" dominant-baseline="middle"
        fill="rgba(255,255,255,0.7)" font-size="11">
        {{a.label}}
      </text>
    </svg>
  `
};
