import { useEffect, useState } from 'react'
import AppLink from '@/shared/ui/AppLink'

const slides = [
  { image: '/carousel/1.jpg', title: '把实践中的问题与答案，沉淀成可复用的知识', label: 'PERSONAL KNOWLEDGE BASE' },
  { image: '/carousel/2.jpg', title: '持续记录，才能让思考留下清晰的轨迹', label: 'NOTES ON BUILDING AND LEARNING' },
]

export default function HomePage() {
  const [activeSlide, setActiveSlide] = useState(0)

  useEffect(() => {
    if (slides.length < 2) return
    const timer = window.setInterval(() => setActiveSlide((current) => (current + 1) % slides.length), 5000)
    return () => window.clearInterval(timer)
  }, [])

  const slide = slides[activeSlide]
  return <main className="home-page">
    <section className="home-cosmic-hero" aria-label="首页轮播图">
      <img className="home-cosmic-hero-image" src={slide.image} alt="瑞克和莫迪的宇宙场景" />
      <div className="home-cosmic-shade" />
      <div className="home-cosmic-content">
        <small>{slide.label}</small><h1>{slide.title}</h1>
      </div>
      <div className="home-carousel-dots" aria-label="轮播图页码">{slides.map((item, index) => <button key={item.image} type="button" className={index === activeSlide ? 'active' : ''} onClick={() => setActiveSlide(index)} aria-label={`查看第 ${index + 1} 张轮播图`} />)}</div>
    </section>
    <section className="home-after-hero">
      <div className="home-shortcuts"><AppLink to="/blog"><b>01</b><div><strong>知识文章</strong><span>KNOWLEDGE ARTICLES</span></div></AppLink><AppLink to="/assistant"><b>02</b><div><strong>智能助手</strong><span>ASK THE ASSISTANT</span></div></AppLink><AppLink to="/diary"><b>03</b><div><strong>随手记录</strong><span>DAILY NOTES</span></div></AppLink><AppLink to="/about"><b>04</b><div><strong>关于本站</strong><span>ABOUT THIS SITE</span></div></AppLink></div>
      <div className="home-section-heading"><small>EXPLORE THE ARCHIVE</small><h2>从一次记录，开始整理自己的知识</h2></div>
      <div className="home-explore-grid"><AppLink to="/blog"><span>01</span><h3>技术沉淀</h3><p>记录开发过程中可复用的方案、踩坑与思考。</p></AppLink><AppLink to="/diary"><span>02</span><h3>项目笔记</h3><p>把零散的过程整理成以后也能读懂的线索。</p></AppLink><AppLink to="/assistant"><span>03</span><h3>知识助手</h3><p>从已有内容中检索答案，让知识保持可访问。</p></AppLink></div>
    </section>
  </main>
}
