from transformers import pipeline

# Load DistilGPT-2 (lighter version)
generator = pipeline("text-generation", model="distilgpt2")

response = generator("Once upon a time", max_length=50, num_return_sequences=1)
print(response[0]['generated_text'])

