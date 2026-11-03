const activities = [
  { title: '社团开放日：和新成员见面', type: '文章发布', time: '今天 09:42', color: 'coral' },
  { title: '「摄影」分类新增 3 篇文章', type: '内容更新', time: '昨天 16:28', color: 'blue' },
  { title: '林默赞了你的文章', type: '互动消息', time: '昨天 11:06', color: 'green' },
  { title: '知识库空间已创建', type: '系统动态', time: '周一 18:20', color: 'purple' },
]

export default function Dashboard() {
  return <div className="dashboard">
    <section className="welcome-panel"><div><span className="eyebrow">周三，2026 年 8 月 10 日</span><h2>早上好，开始整理今天的内容吧。</h2><p>这里是 Renai 社团的内容工作台，记录灵感、发布文章，也和成员保持连接。</p></div><div className="welcome-orbit"><span>R</span></div></section>
    <div className="stat-grid"><div className="stat-card"><span>已发布文章</span><strong>128</strong><small className="positive">↑ 12.5% <em>较上月</em></small></div><div className="stat-card"><span>本月阅读量</span><strong>24,680</strong><small className="positive">↑ 8.2% <em>较上月</em></small></div><div className="stat-card"><span>待处理互动</span><strong>16</strong><small className="warning">需要你的关注</small></div></div>
    <div className="dashboard-columns"><section className="content-panel"><div className="panel-heading"><div><span className="eyebrow">RECENT ACTIVITY</span><h3>最近动态</h3></div><button className="text-action" type="button">查看全部 →</button></div><div className="activity-list">{activities.map((item) => <div className="activity-item" key={item.title}><span className={`activity-mark ${item.color}`} /><div><b>{item.title}</b><small>{item.type} · {item.time}</small></div><span className="activity-more">···</span></div>)}</div></section><section className="content-panel quick-panel"><div className="panel-heading"><div><span className="eyebrow">QUICK ACTIONS</span><h3>快捷操作</h3></div></div><button className="quick-action primary" type="button"><span>✎</span><div><b>写一篇新文章</b><small>分享社团的鲜活内容</small></div><i>→</i></button><button className="quick-action" type="button"><span>☷</span><div><b>整理内容分类</b><small>让文章更容易被发现</small></div><i>→</i></button><button className="quick-action" type="button"><span>◌</span><div><b>进入知识库</b><small>即将上线的内容空间</small></div><i>→</i></button></section></div>
  </div>
}
