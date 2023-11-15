package models

/*

CREATE TABLE "public"."t_proxys" (
  "id" int4 NOT NULL DEFAULT nextval('t_proxys_id_seq'::regclass),
  "created_on" int4 DEFAULT 0,
  "updated_on" int4 DEFAULT 0,
  "deleted_on" int4 DEFAULT 0,
  "is_deleted" bool NOT NULL DEFAULT false,
  "ip" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "host" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "port" int4 NOT NULL,
  "username" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "password" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "status" int2 NOT NULL DEFAULT 101,
  "is_static" bool NOT NULL DEFAULT false,
  "is_traffic" bool NOT NULL DEFAULT false,
  "network" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_datacenter" bool NOT NULL DEFAULT false,
  "city_id" int4 NOT NULL,
  "country_id" int4 NOT NULL,
  "state_id" int4 NOT NULL,
  "up_stream_id" int4 NOT NULL,
  CONSTRAINT "t_proxys_pkey" PRIMARY KEY ("id")
)
;
*/

type ProxyIP struct {
	Model
	Host       string `json:"host"`
	Ip         string `json:"ip"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Port       int    `json:"port"`
	Status     int    `json:"status"`
	UpStreamId int    `json:"up_stream_id"`
}

//type ProxyIP struct {
//	Model
//	Username string `json:"username"`
//	Password string `json:"password"`
//	//PmID        int    `json:"pm_id"`
//	//ForwardPort int    `json:"forward_port"`
//	Port int `json:"port"`
//	//IpID        int    `json:"ip_id"`
//	//Online      int    `json:"online"`
//	//Health      int    `json:"health"`
//	Status int `json:"status"`
//}

func (ProxyIP) TableName() string {
	return "t_proxys"
}

// GetProxyIPByID get a single proxy_ip based on id
func GetProxyIPByID(id int) (*ProxyIP, error) {
	var result ProxyIP
	if err := db.Where("id = ? and deleted_on = ?", id, 0).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
