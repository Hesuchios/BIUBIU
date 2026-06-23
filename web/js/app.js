const { createApp } = Vue;

const app = createApp({
  data() {
    return { page: 'top10' };
  }
});

app.component('top-ranking', TopRanking);
app.component('weapon-db', WeaponDB);
app.component('attachment-db', AttachmentDB);
app.component('mod-code-lib', ModCodeLib);
app.component('smart-mod', SmartMod);
app.component('knowledge-base', KnowledgeBase);
app.component('radar-chart', RadarChart);

app.mount('#app');
