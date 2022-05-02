## v0.9.3

```
insert into sys_apis (created_at, updated_at, deleted_at, path, description, api_group, method) VALUES 
(current_timestamp, current_timestamp, null, '/email/botTest', '企业微信测试', 'email', 'GET');
```

## v0.9.4
```
alter table search_result modify path varchar(800) default '' not null;
alter table token drop column limit_times;
alter table token drop column remaining;
alter table token drop column reset_time;
alter table token drop column description;
```
