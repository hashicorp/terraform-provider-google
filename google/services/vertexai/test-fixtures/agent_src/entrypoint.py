import os
import uvicorn
from fastapi import FastAPI, Request
from fastapi.responses import StreamingResponse
from simple_agent import fixed_name_generator

app = FastAPI()

@app.get("/")
@app.get("/healthz")
def health_check():
    return {"status": "ok"}

@app.post("/api/reasoning_engine")
async def reasoning_engine(request: Request):
    return {"status": "ok"}

@app.post("/api/stream_reasoning_engine")
async def stream_reasoning_engine(request: Request):
    try:
        data = await request.json()
    except Exception:
        data = {}
    product = data.get("kwargs", {}).get("product", "") or data.get("product", "")
    return StreamingResponse(
        fixed_name_generator.stream_query(product),
        media_type="application/json"
    )

@app.get("/stream_query")
def stream_query(product: str = ""):
    return StreamingResponse(
        fixed_name_generator.stream_query(product),
        media_type="application/json"
    )

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8080))
    uvicorn.run(app, host="0.0.0.0", port=port)
