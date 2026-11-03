export interface VectorCollection { name: string; count: number }
export interface VectorDocument { chunk_id: string; doc_id: string; chunk_no: number; title: string; source: string; knowledge_type: string; status: number; created_at?: number; updated_at?: number; content: string }
export interface VectorDocumentPage { total: number; page: number; page_size: number; items: VectorDocument[] }
