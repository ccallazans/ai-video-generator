import sys
from ollama import Client

def generate_text(message):
    client = Client(host='ollama:11434')
    response = client.chat(model='orca-mini', messages=[{'role': 'user', 'content': message}])
    print(response['message']['content'], flush=True)
 
if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python main.py <parameter>")
        sys.exit(1)

    message = sys.argv[1]
    generate_text(message)