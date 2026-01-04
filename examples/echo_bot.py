import asyncio
import json
import websockets

# Configuration
# Replace these with your actual Bot Access Token and connection details
# You can get the Access Token from the "Bot Manager" in the web dashboard.
ACCESS_TOKEN = "YOUR_BOT_ACCESS_TOKEN" 
WS_URL = f"wss://incitymega.cn/bot/v1/onebot/v12/ws?access_token={ACCESS_TOKEN}"

async def echo_bot():
    print(f"Connecting to {WS_URL}...")
    try:
        async with websockets.connect(WS_URL) as websocket:
            print("Connected successfully!")
            print("Waiting for events...")

            async for message in websocket:
                event = json.loads(message)
                
                # Check if it's a message event
                if event.get("type") == "message":
                    print(f"Received message: {event}")
                    
                    # Extract necessary information
                    detail_type = event.get("detail_type")
                    user_id = event["data"]["user_id"]
                    msg_content = event["data"]["content"] # This is a list of message segments for OneBot 12
                    
                    # Construct reply action
                    action = {
                        "action": "send_message",
                        "params": {
                            "detail_type": detail_type,
                            "message": [
                                {"type": "text", "data": {"text": "Echo: "}}
                            ] + msg_content, # Append original content
                        },
                        "echo": "echo_reply"
                    }

                    if detail_type == "private":
                        action["params"]["user_id"] = user_id
                    elif detail_type == "group":
                        action["params"]["group_id"] = event["data"]["group_id"]

                    # Send reply
                    print(f"Sending reply: {json.dumps(action)}")
                    await websocket.send(json.dumps(action))

                # Handle other events (e.g., heartbeats, lifecycle) if needed
                elif event.get("type") == "meta":
                     print(f"Meta event: {event.get('detail_type')}")

    except Exception as e:
        print(f"Connection failed: {e}")
        print("Please check if:")
        print("1. The backend server is running on localhost:8080")
        print("2. The ACCESS_TOKEN is correct")
        print("3. You have pip installed websockets (pip install websockets)")

if __name__ == "__main__":
    # Ensure you have 'websockets' installed:
    # pip install websockets
    try:
        asyncio.run(echo_bot())
    except KeyboardInterrupt:
        print("\nBot stopped.")
