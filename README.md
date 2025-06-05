# Watchtower Metrics Companion

Watchtower强制要求`/v1/metrics`接口使用Bearer, `Token`本身用于防止`api-update`误操作。
因此实现/v1/metrics接口， 透传请求， 并增加HTTP Header Authorization.
