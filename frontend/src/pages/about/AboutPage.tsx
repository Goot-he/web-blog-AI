import AppLink from '@/shared/ui/AppLink'
import './AboutPage.css'

const features = [
  { title: '文章与记录', description: '整理技术实践、项目经验和生活随笔，并通过分类方便查找。' },
  { title: '评论与互动', description: '围绕文章交流观点，支持评论、回复和点赞。' },
  { title: '个人空间', description: '提供登录、个人资料、头像上传和内容管理。' },
  { title: '智能检索', description: '将站内内容接入知识库，通过助手完成检索和问答。' },
]

const projectParts = [
  { title: 'Web 前端', description: '负责文章阅读、评论互动、个人中心和内容管理页面。' },
  { title: 'Gateway 网关', description: '统一处理 HTTP 请求、登录鉴权、文件访问和服务调用。' },
  { title: 'Core 核心服务', description: '承载用户、文章、评论等核心数据和业务处理。' },
  { title: 'Agent 智能服务', description: '负责知识库接入、内容检索和站内问答能力。' },
]

const technologies = [
  'React',
  'TypeScript',
  'Go',
  'Gin',
  'gRPC',
  'MySQL',
  'Redis',
  'FastAPI',
  'ChromaDB',
  'GLM-4-Flash',
]

const improvements = [
  { number: '01', title: '内容体验', description: '继续完善分类、归档、搜索和文章阅读细节。' },
  { number: '02', title: '个人空间', description: '补全资料维护、头像更新和账号相关体验。' },
  { number: '03', title: '智能检索', description: '优化知识库构建、回答质量和内容引用。' },
]

export default function AboutPage() {
  return (
    <main className="inner-page about-page">
      <section className="about-profile">
        <img src="/favicon.svg" alt="Renai 标识" />
        <div>
          <span className="about-eyebrow">ABOUT RENAI</span>
          <h1>关于 Renai</h1>
          <p>记录开发、生活与灵感，也尝试让写下来的知识更容易再次被找到。</p>
          <div className="about-status">
            <span>小型博客项目</span>
            <span>持续开发中</span>
          </div>
        </div>
      </section>

      <section className="about-block">
        <h2>关于这个项目</h2>
        <div className="about-copy">
          <p>
            Renai 是一个用于记录和分享的小型博客。这里既可以发布技术文章与生活记录，也包含评论互动、个人中心和内容管理等基础功能。
          </p>
          <p>
            项目还接入了站内助手和向量知识库，希望让内容不只被阅读一次，而是可以在之后通过搜索和提问重新使用。现在的功能仍在逐步完善，页面和交互也会根据实际使用继续调整。
          </p>
        </div>
      </section>

      <section className="about-block">
        <h2>目前包含</h2>
        <div className="about-feature-grid">
          {features.map((feature) => (
            <article key={feature.title}>
              <h3>{feature.title}</h3>
              <p>{feature.description}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="about-block">
        <h2>项目组成</h2>
        <p className="about-section-note">各部分保持清晰分工，通过网关和 RPC 连接成完整的博客服务。</p>
        <div className="about-part-list">
          {projectParts.map((part) => (
            <article key={part.title}>
              <h3>{part.title}</h3>
              <p>{part.description}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="about-block">
        <h2>项目技术</h2>
        <p className="about-section-note">只展示当前项目中具有代表性的核心技术。</p>
        <div className="about-technology-panel">
          {technologies.map((technology) => (
            <span key={technology}>{technology}</span>
          ))}
        </div>
      </section>

      <section className="about-block">
        <h2>持续完善</h2>
        <div className="about-improvement-grid">
          {improvements.map((item) => (
            <article key={item.number}>
              <span>{item.number}</span>
              <h3>{item.title}</h3>
              <p>{item.description}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="about-footer-block">
        <div>
          <h2>继续看看</h2>
          <p>从一篇文章开始阅读，或者直接向站内助手提问。</p>
        </div>
        <div className="about-actions">
          <AppLink to="/blog">浏览文章</AppLink>
          <AppLink to="/assistant">站内助手</AppLink>
        </div>
      </section>
    </main>
  )
}
