FROM python:3.11-slim

WORKDIR /app

COPY simple_agent.py .
COPY entrypoint.py .
COPY requirements.txt .
RUN pip install -r requirements.txt

CMD ["python", "entrypoint.py"]