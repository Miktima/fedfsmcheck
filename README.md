Script for checking and downloading lists of the federal sites. Notify by emails and telegram when new information occur.
Example of emails.json
```
[
    {
        "list": "ULadd",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-add"
    },
    {
        "list": "FLadd",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-add"
    },
    {
        "list": "ULdel",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-del"
    },
    {
        "list": "FLdel",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-del"
    },
    {
        "list": "Minjust",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://minjust.gov.ru/ru/documents/7756/"
    },
    {
        "list": "Spimex",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://spimex.com/stock_information/market_review/"
    },
    {
        "list": "Acra",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://www.acra-ratings.ru/press-releases/"
    },
    {
        "list": "Mintrans",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://mintrans.gov.ru/press-center/news"
    }
]

```
Example of botkey.json
```
{
    "APIkey": "TELEGRAM_BOT_TOKEN"
}
```