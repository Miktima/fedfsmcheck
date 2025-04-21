Script for checking and downloading lists of the federal sites. Notify by emails and telegram when new information occur.
Docker version
Gocron (https://github.com/go-co-op/gocron) is used for jib scheduling
Example of emails.json
```
[
    {
        "list": "ULadd",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-add"
    },
    {
        "list": "FLadd",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-add"
    },
    {
        "list": "ULdel",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-del"
    },
    {
        "list": "FLdel",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://fedsfm.ru/documents/terrorists-catalog-portal-del"
    },
    {
        "list": "Minjust",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://minjust.gov.ru/ru/documents/7756/"
    },
    {
        "list": "Spimex",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://spimex.com/stock_information/market_review/"
    },
    {
        "list": "Acra",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://www.acra-ratings.ru/press-releases/"
    },
    {
        "list": "Mintrans",
        "cron": "*/10 * * * *",
        "emails": ["email1", "email2"],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://mintrans.gov.ru/press-center/news"
    },
    {
        "list": "Test",
        "cron": "0 */1 * * *",
        "emails": [],
        "chats": ["chat_id1", "chat_id2"],
        "url": "https://test.com"
    }    

]

The last job is used form test of live of the docker

```
Example of botkey.json
```
{
    "APIkey": "TELEGRAM_BOT_TOKEN"
}
```

Data are saved in ./data folder