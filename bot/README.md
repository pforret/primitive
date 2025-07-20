what # Primitive Twitter Bot

A Twitter bot that automatically generates primitive art using geometric shapes and responds to user mentions with custom primitive art generation.

## Overview

This bot provides two main functions:
1. **Automatic Posts**: Fetches interesting photos from Flickr and generates primitive art every 30 minutes
2. **Interactive Mentions**: Responds to user mentions by processing their attached images with custom parameters

## Features

### 🤖 Automated Content Generation
- Fetches random "interesting" photos from Flickr API
- Generates primitive art with randomized parameters
- Posts results to Twitter every 30 minutes
- Links back to original Flickr photos

### 💬 Interactive Mention Handling
- Responds to Twitter mentions containing images
- Parses user requests for specific shapes and counts
- Generates custom primitive art based on user preferences
- Rate-limited to prevent spam (5-minute cooldown per user)

### 🎨 Shape Mode Support
The bot supports all primitive algorithm modes:
- **Triangles** (`m=1`)
- **Rectangles** (`m=2`) 
- **Ellipses** (`m=3`)
- **Circles** (`m=4`)
- **Rotated Rectangles** (`m=5`)
- **Bezier Curves** (`m=6`)
- **Rotated Ellipses** (`m=7`)
- **Polygons** (`m=8`)
- **Mixed/Combo** (`m=0`)

## Setup

### Prerequisites
1. Python 2.7 (legacy code)
2. Required Python packages:
   ```bash
   pip install requests python-twitter
   ```
3. Compiled `primitive` binary in PATH
4. Twitter API credentials
5. Flickr API key

### Configuration
Create a `config.py` file in the bot directory:

```python
# Twitter API Credentials
TWITTER_CONSUMER_KEY = 'your_consumer_key'
TWITTER_CONSUMER_SECRET = 'your_consumer_secret' 
TWITTER_ACCESS_TOKEN_KEY = 'your_access_token'
TWITTER_ACCESS_TOKEN_SECRET = 'your_access_token_secret'

# Flickr API
FLICKR_API_KEY = 'your_flickr_api_key'

# File Paths
INPUT_FOLDER = '/path/to/input/images'
OUTPUT_FOLDER = '/path/to/output/images'
```

### Directory Setup
```bash
mkdir -p /path/to/input/images
mkdir -p /path/to/output/images
```

## Usage

### Running the Bot
```bash
cd bot/
python main.py
```

The bot will run continuously, checking for mentions every 65 seconds and generating automatic posts every 30 minutes.

### User Interaction Examples

**Basic mention:**
```
@your_bot_handle [image attached]
```
Bot responds with randomized primitive art.

**Specific shape request:**
```
@your_bot_handle 100 triangles [image attached]
```
Bot responds: `@username 100 triangles.` + generated image

**Shape type keywords:**
- `triangles` → Triangle mode
- `rectangles` → Rectangle mode  
- `ellipses` → Ellipse mode
- `circles` → Circle mode
- `beziers` → Bezier curve mode
- `polygons` → Polygon mode

**Number parsing:**
Any number in the mention (1-500) will be used as the shape count.

## Technical Details

### Main Components

#### `Config` Class
Handles parameter parsing and validation:
```python
config = Config()
config.randomize()        # Set random parameters
config.parse(tweet_text)  # Parse user request
config.validate()         # Ensure valid ranges
```

#### Core Functions
- `generate()` - Automatic Flickr photo processing
- `handle_mentions()` - Twitter mention processing
- `handle_mention(status)` - Individual mention handler
- `primitive(**kwargs)` - Wrapper for primitive binary

#### Safety Features
- **Rate limiting**: 5-minute cooldown per user
- **Validation**: Requires exactly one mention and one image
- **Timestamp checking**: Ignores old mentions
- **Error handling**: Continues operation despite failures

### File Structure
```
bot/
├── main.py           # Main bot code
├── config.py         # API credentials (create this)
├── requirements.txt  # Python dependencies
└── README.md         # This file
```

### Parameters

#### Default Randomization
```python
self.m = random.choice([1, 5, 6, 7])  # Shape modes
self.n = random.randint(10, 50) * 10   # Shape count (100-500)
self.a = 128                           # Alpha value
self.r = 300                           # Image resize
self.s = 1200                          # Image size
```

#### Special Cases
- **Bezier mode** (`m=6`): Fixed to 100 shapes with alpha=0, rep=19
- **Shape count**: Clamped to 1-500 range
- **Mode validation**: Clamped to 0-8 range

## API Requirements

### Twitter API
Requires a Twitter Developer account and app with:
- Consumer Key/Secret
- Access Token/Secret
- Read and Write permissions

### Flickr API
Requires Flickr API key for accessing the "interestingness" feed:
- Free tier sufficient for bot usage
- Used to fetch random interesting photos

## Rate Limits

- **Automatic posts**: Every 30 minutes (48 posts/day max)
- **Mention checking**: Every 65 seconds
- **User rate limit**: 5-minute cooldown between mentions per user
- **Twitter API limits**: Respects standard Twitter API rate limits

## Deployment Notes

### Production Considerations
1. **Process management**: Use `systemd`, `supervisor`, or similar
2. **Logging**: Redirect output to log files
3. **Monitoring**: Track bot health and API usage
4. **Storage**: Ensure adequate disk space for image processing
5. **Cleanup**: Implement periodic cleanup of processed images

### Example systemd service
```ini
[Unit]
Description=Primitive Twitter Bot
After=network.target

[Service]
Type=simple
User=bot
WorkingDirectory=/path/to/primitive/bot
ExecStart=/usr/bin/python main.py
Restart=always

[Install]
WantedBy=multi-user.target
```

## Troubleshooting

### Common Issues
1. **Missing config**: Ensure `config.py` exists with valid API credentials
2. **Primitive binary**: Verify `primitive` command is in PATH
3. **Directory permissions**: Ensure bot can read/write to input/output folders
4. **API limits**: Monitor Twitter/Flickr API usage
5. **Image processing failures**: Check image format compatibility

### Debug Mode
Add debug prints or modify rate limits for testing:
```python
RATE = 60        # Post every minute for testing
MENTION_RATE = 10  # Check mentions every 10 seconds
```

## Contributing

To modify or extend the bot:
1. **Add new shape modes**: Update `MODE_NAMES` array
2. **Modify parameters**: Edit `Config.randomize()` method
3. **Change posting frequency**: Modify `RATE` constant
4. **Add new features**: Extend mention parsing in `Config.parse()`

## License

This bot is part of the primitive project. See the main project LICENSE for details.