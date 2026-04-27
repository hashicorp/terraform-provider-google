class MetadataAgent:

  def query(self):
    import requests
    url = "http://metadata.google.internal/computeMetadata/v1/project/numeric-project-id"
    try:
        response = requests.get(url, headers={"Metadata-Flavor": "Google"})
        response.raise_for_status()
        return f"service-{response.text}@serverless-robot-prod.iam.gserviceaccount.com"
    except Exception:
        return None

root_agent = MetadataAgent()