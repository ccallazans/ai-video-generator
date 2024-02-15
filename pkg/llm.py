import sys
import ollama

def generate_text(message):
    response = ollama.chat(model='orca-mini', messages=[{'role': 'user', 'content': message}])
    print(response['message']['content'], flush=True)
 
if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python main.py <parameter>")
        sys.exit(1)

    message = sys.argv[1]
    generate_text(message)