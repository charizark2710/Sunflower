#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>

// Configuration struct to hold module settings (like DB connection strings)
typedef struct
{
    ngx_flag_t enable;
    ngx_str_t db_connection_string;
} ngx_http_hmac_db_loc_conf_t;

static ngx_int_t ngx_http_hmac_db_handler(ngx_http_request_t *r);
static void *ngx_http_hmac_db_create_loc_conf(ngx_conf_t *cf);
static char *ngx_http_hmac_db_merge_loc_conf(ngx_conf_t *cf, void *parent, void *child);
static char *ngx_http_hmac_db_enable(ngx_conf_t *cf, ngx_command_t *cmd, void *conf);

// Directives: This defines how your module is invoked in nginx.conf
static ngx_command_t ngx_http_hmac_db_commands[] = {
    {ngx_string("hmac_db_verify"),
     NGX_HTTP_LOC_CONF | NGX_CONF_FLAG,
     ngx_http_hmac_db_enable,
     NGX_HTTP_LOC_CONF_OFFSET,
     offsetof(ngx_http_hmac_db_loc_conf_t, enable),
     NULL},

    ngx_null_command};

// Module Context
static ngx_http_module_t ngx_http_hmac_db_module_ctx = {
    NULL,                             /* preconfiguration */
    NULL,                             /* postconfiguration */
    NULL,                             /* create main configuration */
    NULL,                             /* init main configuration */
    NULL,                             /* create server configuration */
    NULL,                             /* merge server configuration */
    ngx_http_hmac_db_create_loc_conf, /* create location configuration */
    ngx_http_hmac_db_merge_loc_conf   /* merge location configuration */
};

// Module Definition
ngx_module_t ngx_http_hmac_db_module = {
    NGX_MODULE_V1,
    &ngx_http_hmac_db_module_ctx, /* module context */
    ngx_http_hmac_db_commands,    /* module directives */
    NGX_HTTP_MODULE,              /* module type */
    NULL,                         /* init master */
    NULL,                         /* init module */
    NULL,                         /* init process */
    NULL,                         /* init thread */
    NULL,                         /* exit thread */
    NULL,                         /* exit process */
    NULL,                         /* exit master */
    NGX_MODULE_V1_PADDING};

// The Actual Handler (Where the magic—and database lookup—happens)
static ngx_int_t
ngx_http_hmac_db_handler(ngx_http_request_t *r)
{
    // 1. Read incoming HMAC from headers/query string
    // 2. Query Memcached/Postgres (Must be non-blocking!)
    // 3. Hash the retrieved data
    // 4. Compare hashes

    // If verification fails:
    // return NGX_HTTP_FORBIDDEN;

    // If verification passes, decline handling to let Nginx pass the request to the next phase (or upstream)
    return NGX_DECLINED;
}

// Helper to register the handler to the ACCESS phase
static char *
ngx_http_hmac_db_enable(ngx_conf_t *cf, ngx_command_t *cmd, void *conf)
{
    ngx_http_core_loc_conf_t *clcf;

    clcf = ngx_http_conf_get_module_loc_conf(cf, ngx_http_core_module);
    clcf->handler = ngx_http_hmac_db_handler;

    return NGX_CONF_OK;
}

// Memory allocation for configuration
static void *
ngx_http_hmac_db_create_loc_conf(ngx_conf_t *cf)
{
    ngx_http_hmac_db_loc_conf_t *conf;

    conf = ngx_pcalloc(cf->pool, sizeof(ngx_http_hmac_db_loc_conf_t));
    if (conf == NULL)
    {
        return NULL;
    }
    conf->enable = NGX_CONF_UNSET;
    return conf;
}

static char *
ngx_http_hmac_db_merge_loc_conf(ngx_conf_t *cf, void *parent, void *child)
{
    ngx_http_hmac_db_loc_conf_t *prev = parent;
    ngx_http_hmac_db_loc_conf_t *conf = child;

    ngx_conf_merge_value(conf->enable, prev->enable, 0);
    return NGX_CONF_OK;
}