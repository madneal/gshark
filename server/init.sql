create table casbin_rule
(
	p_type varchar(100) null,
	v0 varchar(100) null,
	v1 varchar(100) null,
	v2 varchar(100) null,
	v3 varchar(100) null,
	v4 varchar(100) null,
	v5 varchar(100) null
);

create table exa_customers
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	customer_name varchar(191) null comment '客户名',
	customer_phone_data varchar(191) null comment '客户手机号',
	sys_user_id bigint unsigned null comment '管理ID',
	sys_user_authority_id varchar(191) null comment '管理角色ID'
);

create index idx_exa_customers_deleted_at
	on exa_customers (deleted_at);

create table exa_file_chunks
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	exa_file_id bigint unsigned null,
	file_chunk_number bigint null,
	file_chunk_path varchar(191) null
);

create index idx_exa_file_chunks_deleted_at
	on exa_file_chunks (deleted_at);

create table exa_file_upload_and_downloads
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	name varchar(191) null comment '文件名',
	url varchar(191) null comment '文件地址',
	tag varchar(191) null comment '文件标签',
	`key` varchar(191) null comment '编号'
);

create index idx_exa_file_upload_and_downloads_deleted_at
	on exa_file_upload_and_downloads (deleted_at);

create table exa_files
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	file_name varchar(191) null,
	file_md5 varchar(191) null,
	file_path varchar(191) null,
	chunk_total bigint null,
	is_finish tinyint(1) null
);

create index idx_exa_files_deleted_at
	on exa_files (deleted_at);

create table exa_simple_uploaders
(
	chunk_number varchar(191) null comment '当前切片标记',
	current_chunk_size varchar(191) null comment '当前切片容量',
	current_chunk_path varchar(191) null comment '切片本地路径',
	total_size varchar(191) null comment '总容量',
	identifier varchar(191) null comment '文件标识（md5）',
	filename varchar(191) null comment '文件名',
	total_chunks varchar(191) null comment '切片总数',
	is_done tinyint(1) null comment '是否上传完成',
	file_path varchar(191) null comment '文件本地路径'
);

create table filters
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	extension varchar(100) null,
	whiteExts varchar(100) null
);

create index idx_filters_deleted_at
	on filters (deleted_at);

create table jwt_blacklists
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	jwt text null comment 'jwt'
);

create index idx_jwt_blacklists_deleted_at
	on jwt_blacklists (deleted_at);

create table repo
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	project_id bigint null,
	type varchar(10) null,
	`desc` varchar(300) null,
	url varchar(300) null,
	path varchar(100) null,
	last_activity_at datetime null,
	status bigint null
);

create index idx_repo_deleted_at
	on repo (deleted_at);

create table rule
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	type varchar(20) null,
	content varchar(100) null,
	name varchar(100) null,
	`desc` varchar(300) null,
	status tinyint null
);

create index idx_rule_deleted_at
	on rule (deleted_at);

create table search_result
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	repo varchar(200) null,
	repository varchar(200) null,
	matches text null,
	keyword varchar(100) null,
	path varchar(500) null,
	url varchar(500) null,
	textmatch_md5 varchar(100) null,
	status tinyint null,
	text_matches_json json null
);

create index idx_search_result_deleted_at
	on search_result (deleted_at);

create table subdomain
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	subdomain varchar(100) null,
	domain varchar(100) null,
	status tinyint null
);

create index idx_subdomain_deleted_at
	on subdomain (deleted_at);

create table sys_apis
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	path varchar(191) null comment 'api路径',
	description varchar(191) null comment 'api中文描述',
	api_group varchar(191) null comment 'api组',
	method varchar(191) default 'POST' null
);

create index idx_sys_apis_deleted_at
	on sys_apis (deleted_at);

create table sys_authorities
(
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	authority_id varchar(90) not null comment '角色ID',
	authority_name varchar(191) null comment '角色名',
	parent_id varchar(191) null comment '父角色ID',
	default_router varchar(191) default 'dashboard' null comment '默认菜单',
	constraint authority_id
		unique (authority_id)
);

alter table sys_authorities
	add primary key (authority_id);

create table sys_authority_menus
(
	sys_base_menu_id bigint unsigned not null,
	sys_authority_authority_id varchar(90) not null comment '角色ID',
	primary key (sys_base_menu_id, sys_authority_authority_id)
);

create table sys_base_menu_parameters
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	sys_base_menu_id bigint unsigned null,
	type varchar(191) null comment '地址栏携带参数为params还是query',
	`key` varchar(191) null comment '地址栏携带参数的key',
	value varchar(191) null comment '地址栏携带参数的值'
);

create index idx_sys_base_menu_parameters_deleted_at
	on sys_base_menu_parameters (deleted_at);

create table sys_base_menus
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	menu_level bigint unsigned null,
	parent_id varchar(191) null comment '父菜单ID',
	path varchar(191) null comment '路由path',
	name varchar(191) null comment '路由name',
	hidden tinyint(1) null comment '是否在列表隐藏',
	component varchar(191) null comment '对应前端文件路径',
	sort bigint null comment '排序标记',
	keep_alive tinyint(1) null comment '附加属性',
	default_menu tinyint(1) null comment '附加属性',
	title varchar(191) null comment '附加属性',
	icon varchar(191) null comment '附加属性',
	close_tab tinyint(1) null comment '附加属性'
);

create index idx_sys_base_menus_deleted_at
	on sys_base_menus (deleted_at);

create table sys_data_authority_id
(
	sys_authority_authority_id varchar(90) not null comment '角色ID',
	data_authority_id_authority_id varchar(90) not null comment '角色ID',
	primary key (sys_authority_authority_id, data_authority_id_authority_id)
);

create table sys_dictionaries
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	name varchar(191) null comment '字典名（中）',
	type varchar(191) null comment '字典名（英）',
	status tinyint(1) null comment '状态',
	`desc` varchar(191) null comment '描述'
);

create index idx_sys_dictionaries_deleted_at
	on sys_dictionaries (deleted_at);

create table sys_dictionary_details
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	label varchar(191) null comment '展示值',
	value bigint null comment '字典值',
	status tinyint(1) null comment '启用状态',
	sort bigint null comment '排序标记',
	sys_dictionary_id bigint unsigned null comment '关联标记'
);

create index idx_sys_dictionary_details_deleted_at
	on sys_dictionary_details (deleted_at);

create table sys_operation_records
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	ip varchar(191) null comment '请求ip',
	method varchar(191) null comment '请求方法',
	path varchar(191) null comment '请求路径',
	status bigint null comment '请求状态',
	latency bigint null comment '延迟',
	agent varchar(191) null comment '代理',
	error_message varchar(191) null comment '错误信息',
	body longtext null comment '请求Body',
	resp longtext null comment '响应Body',
	user_id bigint unsigned null comment '用户id'
);

create index idx_sys_operation_records_deleted_at
	on sys_operation_records (deleted_at);

create table sys_users
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	uuid varchar(191) null comment '用户UUID',
	username varchar(191) null comment '用户登录名',
	password varchar(191) null comment '用户登录密码',
	nick_name varchar(191) default '系统用户' null comment '用户昵称',
	header_img varchar(191) default 'http://qmplusimg.henrongyi.top/head.png' null comment '用户头像',
	authority_id varchar(90) default '888' null comment '用户角色ID'
);

create index idx_sys_users_deleted_at
	on sys_users (deleted_at);

create table token
(
	id bigint unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	type varchar(10) null,
	content varchar(100) null,
	description varchar(100) null,
	limit_times int(5) null,
	remaining int(5) null,
	reset_time datetime null
);

create index idx_token_deleted_at
	on token (deleted_at);

