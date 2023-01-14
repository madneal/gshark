## v1.1.0

```
alter table search_result drop column textmatch_md5;

alter table filters drop column extension;
alter table filters drop column whiteExts;
alter table filters drop column keywords;
alter table filters add filter_type varchar(20) default 'black' null;
alter table filters add filter_class varchar(20) default 'extension' null;
alter table filters add content varchar(100) default '' null;

insert into sys_apis (created_at, updated_at, deleted_at, path, description, api_group, method) VALUES
    (current_timestamp, current_timestamp, null, '/searchResult/startSecFilterTask', '开始二次过滤任务', 'searchResult', 'POST');
insert into casbin_rule (p_type, v0, v1, v2) values ('p', 888, '/searchResult/startSecFilterTask', 'POST');

insert into sys_apis (created_at, updated_at, deleted_at, path, description, api_group, method) VALUES
    (current_timestamp, current_timestamp, null, '/searchResult/getTaskStatus', '开始二次过滤任务', 'searchResult', 'GET');
insert into casbin_rule (p_type, v0, v1, v2) values ('p', 888, '/searchResult/getTaskStatus', 'GET');
```

## v1.0.3

```
insert into sys_apis (created_at, updated_at, deleted_at, path, description, api_group, method) VALUES
    (current_timestamp, current_timestamp, null, '/rule/uploadRules', '规则导入', 'rule', 'POST');
insert into casbin_rule (p_type, v0, v1, v2) VALUES ('p', '888', '/rule/uploadRules', 'POST')

alter table rule drop column deleted_at;
```

## v0.9.9

```
alter table rule
    modify type varchar(100) null;
alter table rule
    change type rule_type varchar(200) null;
```

## v0.9.8

```
insert into sys_apis (created_at, updated_at, deleted_at, path, description, api_group, method) VALUES
(current_timestamp, current_timestamp, null, '/rule/batchCreateRule', '批量导入规则', 'rule', 'POST');
insert into casbin_rule (p_type, v0, v1, v2) values ('p', 888, '/rule/batchCreateRule', 'POST');
```

```
insert into sys_apis (created_at, updated_at, deleted_at, path, description, api_group, method) VALUES
    (current_timestamp, current_timestamp, null, '/rule/switchRuleStatus', '切换规则状态', 'rule', 'POST');

insert into casbin_rule (p_type, v0, v1, v2) values ('p', 888, '/rule/switchRuleStatus', 'POST');
```

## v0.9.4
```
alter table search_result modify path varchar(800) default '' not null;
alter table token drop column limit_times;
alter table token drop column remaining;
alter table token drop column reset_time;
alter table token drop column description;
```

## v0.9.3

```
insert into sys_apis (created_at, updated_at, deleted_at, path, description, api_group, method) VALUES 
(current_timestamp, current_timestamp, null, '/email/botTest', '企业微信测试', 'email', 'GET');
```







