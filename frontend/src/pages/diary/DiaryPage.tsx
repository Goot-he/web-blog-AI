import { useCallback, useEffect, useRef, useState } from 'react'
import * as echarts from 'echarts'
import chinaMap from 'china-geojson/src/geojson/china.json'

type Moment = { province: string; title: string; subtitle: string; photo: string; color: string; capital: [number, number]; position: { top?: string; right?: string; bottom?: string; left?: string } }
type Point = { x: number; y: number }

const moments: Moment[] = [
  { province: '北京市', title: '北京 · 待补充', subtitle: '胡同里的日常片段', photo: 'BEIJING', color: '#b8d8ee', capital: [116.4074, 39.9042], position: { top: '12%', right: '5%' } },
  { province: '浙江省', title: '西湖 · 待补充', subtitle: '把这一刻留在生活小记', photo: 'HANGZHOU', color: '#b6ded6', capital: [120.1551, 30.2741], position: { bottom: '11%', right: '7%' } },
  { province: '四川省', title: '成都 · 待补充', subtitle: '雨后的城市散步', photo: 'CHENGDU', color: '#f7d1a5', capital: [104.0665, 30.5723], position: { bottom: '14%', left: '5%' } },
  { province: '西藏自治区', title: '西藏 · 待补充', subtitle: '这里的故事还没写下', photo: 'TIBET', color: '#dfc2e7', capital: [91.1322, 29.6604], position: { top: '18%', left: '5%' } },
]

const provinceFiles: Record<string, string> = {
  '北京市': 'bei_jing', '天津市': 'tian_jin', '上海市': 'shang_hai', '重庆市': 'chong_qing', '河北省': 'he_bei', '山西省': 'shan_xi_1', '辽宁省': 'liao_ning', '吉林省': 'ji_lin', '黑龙江省': 'hei_long_jiang', '江苏省': 'jiang_su', '浙江省': 'zhe_jiang', '安徽省': 'an_hui', '福建省': 'fu_jian', '江西省': 'jiang_xi', '山东省': 'shan_dong', '河南省': 'he_nan', '湖北省': 'hu_bei', '湖南省': 'hu_nan', '广东省': 'guang_dong', '海南省': 'hai_nan', '四川省': 'si_chuan', '贵州省': 'gui_zhou', '云南省': 'yun_nan', '陕西省': 'shan_xi_2', '甘肃省': 'gan_su', '青海省': 'qing_hai', '台湾省': 'tai_wan', '内蒙古自治区': 'nei_meng_gu', '广西壮族自治区': 'guang_xi', '西藏自治区': 'xi_zang', '宁夏回族自治区': 'ning_xia', '新疆维吾尔自治区': 'xin_jiang', '香港特别行政区': 'xiang_gang', '澳门特别行政区': 'ao_men',
}
const provinceByID: Record<string, string> = { '11': '北京市', '12': '天津市', '13': '河北省', '14': '山西省', '15': '内蒙古自治区', '21': '辽宁省', '22': '吉林省', '23': '黑龙江省', '31': '上海市', '32': '江苏省', '33': '浙江省', '34': '安徽省', '35': '福建省', '36': '江西省', '37': '山东省', '41': '河南省', '42': '湖北省', '43': '湖南省', '44': '广东省', '45': '广西壮族自治区', '46': '海南省', '50': '重庆市', '51': '四川省', '52': '贵州省', '53': '云南省', '54': '西藏自治区', '61': '陕西省', '62': '甘肃省', '63': '青海省', '64': '宁夏回族自治区', '65': '新疆维吾尔自治区', '71': '台湾省', '81': '香港特别行政区', '82': '澳门特别行政区' }
const provinceModules = import.meta.glob('../../../node_modules/china-geojson/src/geojson/*_geo.json', { eager: true, import: 'default' })
const palette = ['#9edbd8', '#f9d69c', '#d9b7e2', '#a9dbf4', '#f7bbcc', '#d8e6bc', '#f1c4d8']

export default function DiaryPage() {
  const pageRef = useRef<HTMLElement>(null)
  const mapRef = useRef<HTMLDivElement>(null)
  const cardRefs = useRef<Record<string, HTMLElement | null>>({})
  const chartRef = useRef<echarts.ECharts | null>(null)
  const dragRef = useRef<{ province: string; offsetX: number; offsetY: number } | null>(null)
  const lastProvinceClickRef = useRef<{ province: string; at: number } | null>(null)
  const [selected, setSelected] = useState(moments[3])
  const [detailProvince, setDetailProvince] = useState<string | null>(null)
  const [capitalPoints, setCapitalPoints] = useState<Record<string, Point>>({})
  const [cardPoints, setCardPoints] = useState<Record<string, Point>>({})

  const updateConnectionPoints = useCallback(() => {
    const chart = chartRef.current
    const page = pageRef.current
    const map = mapRef.current
    if (!chart || !page || !map || detailProvince) return
    const pageBox = page.getBoundingClientRect()
    const mapBox = map.getBoundingClientRect()
    const nextCapitals: Record<string, Point> = {}
    const nextCards: Record<string, Point> = {}
    moments.forEach((moment) => {
      const pixel = chart.convertToPixel({ seriesIndex: 0 }, moment.capital) as number[]
      const cardBox = cardRefs.current[moment.province]?.getBoundingClientRect()
      if (pixel && Number.isFinite(pixel[0])) nextCapitals[moment.province] = { x: mapBox.left - pageBox.left + pixel[0], y: mapBox.top - pageBox.top + pixel[1] }
      if (cardBox) nextCards[moment.province] = { x: cardBox.left - pageBox.left + cardBox.width / 2, y: cardBox.top - pageBox.top + cardBox.height / 2 }
    })
    setCapitalPoints(nextCapitals)
    setCardPoints(nextCards)
  }, [detailProvince])

  useEffect(() => {
    if (!mapRef.current) return
    const chart = echarts.init(mapRef.current)
    chartRef.current = chart
    const mapName = detailProvince ? `diary-${detailProvince}` : 'china-diary-fullscreen'
    const file = detailProvince && provinceFiles[detailProvince]
    const mapData = (file ? provinceModules[`../../../node_modules/china-geojson/src/geojson/${file}_geo.json`] : chinaMap) as { features: { properties: { name: string; id?: string } }[] }
    echarts.registerMap(mapName, mapData as never)
    chart.setOption({ tooltip: { show: false }, series: [{ type: 'map', map: mapName, roam: true, selectedMode: 'single', layoutCenter: ['50%', '52%'], layoutSize: detailProvince ? '76%' : '46%', label: { show: Boolean(detailProvince), color: '#4e6874', fontSize: 11 }, itemStyle: { borderColor: '#fff', borderWidth: 2 }, emphasis: { label: { show: true, color: '#315260', fontSize: 12, fontWeight: 700 }, itemStyle: { areaColor: '#ffb0c9', shadowBlur: 14, shadowColor: 'rgba(231, 119, 159, .45)' } }, data: mapData.features.map((feature, index) => ({ name: feature.properties.name, itemStyle: { areaColor: palette[index % palette.length] } })) }] })
    const updatePoints = () => { requestAnimationFrame(updateConnectionPoints) }
    const resize = () => { chart.resize(); updatePoints() }
    const resizeObserver = new ResizeObserver(resize)
    resizeObserver.observe(mapRef.current)
    chart.on('finished', updatePoints)
    chart.on('georoam', updatePoints)
    chart.on('click', (params) => {
      if (detailProvince) return
      const id = mapData.features[params.dataIndex]?.properties.id
      const province = provinceByID[String(id)]
      if (!province) return
      const now = Date.now()
      const lastClick = lastProvinceClickRef.current
      if (lastClick?.province === province && now - lastClick.at < 450) {
        setDetailProvince(province)
        lastProvinceClickRef.current = null
        return
      }
      lastProvinceClickRef.current = { province, at: now }
      const moment = moments.find((item) => item.province === province)
      if (moment) setSelected(moment)
    })
    window.addEventListener('resize', resize)
    return () => { resizeObserver.disconnect(); window.removeEventListener('resize', resize); chart.dispose(); chartRef.current = null }
  }, [detailProvince, updateConnectionPoints])

  useEffect(() => { const id = requestAnimationFrame(updateConnectionPoints); return () => cancelAnimationFrame(id) }, [updateConnectionPoints])

  const startDrag = (event: React.PointerEvent<HTMLElement>, province: string) => {
    const box = event.currentTarget.getBoundingClientRect()
    dragRef.current = { province, offsetX: event.clientX - box.left, offsetY: event.clientY - box.top }
    event.currentTarget.setPointerCapture(event.pointerId)
  }
  const moveDrag = (event: React.PointerEvent<HTMLElement>) => {
    const drag = dragRef.current
    const page = pageRef.current
    if (!drag || !page) return
    const pageBox = page.getBoundingClientRect()
    const card = cardRefs.current[drag.province]
    if (!card) return
    const x = Math.max(12, Math.min(pageBox.width - card.offsetWidth - 12, event.clientX - pageBox.left - drag.offsetX))
    const y = Math.max(12, Math.min(pageBox.height - card.offsetHeight - 12, event.clientY - pageBox.top - drag.offsetY))
    Object.assign(card.style, { left: `${x}px`, top: `${y}px`, right: 'auto', bottom: 'auto' })
    updateConnectionPoints()
  }
  const stopDrag = () => { dragRef.current = null }
  const connectionPath = (from: Point, to: Point) => `M ${from.x} ${from.y} C ${from.x + (to.x - from.x) * .45} ${from.y}, ${to.x - (to.x - from.x) * .35} ${to.y}, ${to.x} ${to.y}`

  return <main ref={pageRef} className="diary-fullscreen" aria-label="生活小记地图">
    <div className="diary-map-hint">拖动地图查看 · 连续点击同一省份查看城市</div>
    {detailProvince && <button className="diary-map-back" type="button" onClick={() => setDetailProvince(null)}>返回全国地图</button>}
    <section className="diary-map-stage"><div ref={mapRef} className="diary-full-map" /></section>
    {!detailProvince && <>
      <svg className="diary-dotted-lines" aria-hidden="true">{cardPoints[selected.province] && capitalPoints[selected.province] && <path d={connectionPath(cardPoints[selected.province], capitalPoints[selected.province])} />}</svg>
      <span className="diary-capital-dot" style={{ left: capitalPoints[selected.province]?.x, top: capitalPoints[selected.province]?.y }} aria-hidden="true" />
      {moments.map((moment) => <article ref={(node) => { cardRefs.current[moment.province] = node }} key={moment.province} className={`diary-polaroid${selected.province === moment.province ? ' active' : ''}`} style={{ ...moment.position, display: selected.province === moment.province ? undefined : 'none' }} onPointerDown={(event) => startDrag(event, moment.province)} onPointerMove={moveDrag} onPointerUp={stopDrag} onPointerCancel={stopDrag} onClick={() => setSelected(moment)}>
        <div className="diary-photo-art" style={{ backgroundColor: moment.color }}><span>{moment.photo}</span></div><b>{moment.title.split('·')[0]}</b><small>{moment.subtitle}</small>
      </article>)}
    </>}
  </main>
}
