package model

const VsIPCountryInfoTableName = "vs_ip_country_info"

// 虚拟服务IP国家信息表
type VsIPCountryInfo struct {
	ID                  uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	CountryName         string `gorm:"column:country_name;type:varchar(100);NOT NULL" json:"country_name"`                                    // 国家名称
	LongName            string `gorm:"column:long_name;type:varchar(255)" json:"long_name"`                                                   // 国家英文全称
	Code                string `gorm:"column:code;type:char(20)" json:"code"`                                                                 // 国家代码
	IsDynamic           int    `gorm:"column:is_dynamic;type:tinyint(1);default:1;NOT NULL" json:"is_dynamic"`                                // 是否支持动态ip 0否 1是
	IproyalStaticID     int    `gorm:"column:iproyal_static_id;type:int(11)" json:"iproyal_static_id"`                                        // iproyal静态住宅 id
	IproyalDatacenterID int    `gorm:"column:iproyal_datacenter_id;type:int(11)" json:"iproyal_datacenter_id"`                                // iproyal静态机房 id
	LineV4              int    `gorm:"column:line_v4;type:tinyint(1);NOT NULL" json:"line_v4"`                                                // 是否支持v4 0否 1是
	LineV6              int    `gorm:"column:line_v6;type:tinyint(1);NOT NULL" json:"line_v6"`                                                // 是否支持v6 0否 1是
	Proxy6V4            int    `gorm:"column:proxy6_v4;type:tinyint(1);default:0;NOT NULL" json:"proxy6_v4"`                                  // 是否支持v4 0否 1是
	Proxy6ShareV4       int    `gorm:"column:proxy6_share_v4;type:tinyint(1);default:0;NOT NULL" json:"proxy6_share_v4"`                      // 是否支持共享v4 0否 1是
	Proxy6V6            int    `gorm:"column:proxy6_v6;type:tinyint(1);default:0;NOT NULL" json:"proxy6_v6"`                                  // 是否支持v6 0否 1是
	ExclusiveV4         int    `gorm:"column:exclusive_v4;type:tinyint(1);default:0;NOT NULL" json:"exclusive_v4"`                            // 是否支持独享v4 0否 1是
	ExclusiveV6         int    `gorm:"column:exclusive_v6;type:tinyint(1);default:0;NOT NULL" json:"exclusive_v6"`                            // 是否支持独享v6 0否 1是
	Vkontakte           int    `gorm:"column:vkontakte;type:tinyint(1);default:0;NOT NULL" json:"vkontakte"`                                  // 是否支持VK 0否 1是
	Instagram           int    `gorm:"column:instagram;type:tinyint(1);default:0;NOT NULL" json:"instagram"`                                  // 是否支持ins 0否 1是
	Facebook            int    `gorm:"column:facebook;type:tinyint(1);default:0;NOT NULL" json:"facebook"`                                    // 是否支持fb 0否 1是
	Twitter             int    `gorm:"column:twitter;type:tinyint(1);default:0;NOT NULL" json:"twitter"`                                      // 是否支持twitter 0否 1是
	Steam               int    `gorm:"column:steam;type:tinyint(1);default:0;NOT NULL" json:"steam"`                                          // 是否支持steam 0否 1是
	Odnoklassniki       int    `gorm:"column:odnoklassniki;type:tinyint(1);default:0;NOT NULL" json:"odnoklassniki"`                          // 是否支持odnoklassniki 0否 1是
	Gmail               int    `gorm:"column:gmail;type:tinyint(1);default:0;NOT NULL" json:"gmail"`                                          // 是否支持gmail 0否 1是
	Youtube             int    `gorm:"column:youtube;type:tinyint(1);default:0;NOT NULL" json:"youtube"`                                      // 是否支持油管 0否 1是
	Twitch              int    `gorm:"column:twitch;type:tinyint(1);default:0;NOT NULL" json:"twitch"`                                        // 是否支持twich 0否 1是
	Pinterest           int    `gorm:"column:pinterest;type:tinyint(1);default:0;NOT NULL" json:"pinterest"`                                  // 是否支持pinterset 0否 1是
	Wot                 int    `gorm:"column:wot;type:tinyint(1);default:0;NOT NULL" json:"wot"`                                              // 是否支持wot 0否 1是
	Lineage             int    `gorm:"column:lineage;type:tinyint(1);default:0;NOT NULL" json:"lineage"`                                      // 是否支持血统 0否 1是
	Amazon              int    `gorm:"column:amazon;type:tinyint(1);default:0;NOT NULL" json:"amazon"`                                        // 是否支持亚马逊 0否 1是
	Ebay                int    `gorm:"column:ebay;type:tinyint(1);default:0;NOT NULL" json:"ebay"`                                            // 是否支持ebay 0否 1是
	Wow                 int    `gorm:"column:wow;type:tinyint(1);default:0;NOT NULL" json:"wow"`                                              // 是否支持魔兽世界 0否 1是
	OtherSocialNetworks uint   `gorm:"column:other-social-networks;type:tinyint(1) unsigned;default:0;NOT NULL" json:"other-social-networks"` // 是否支持其他网站 0否 1是
	Fortnite            int    `gorm:"column:fortnite;type:tinyint(1);default:0;NOT NULL" json:"fortnite"`                                    // 是否支持堡垒之夜 0否 1是
	Bookmaker           int    `gorm:"column:bookmaker;type:tinyint(1);default:0;NOT NULL" json:"bookmaker"`                                  // 是否支持博彩 0否 1是
	OtherGames          int    `gorm:"column:other-games;type:tinyint(1);default:0;NOT NULL" json:"other-games"`                              // 是否支持其他游戏 0否 1是
	Tiktok              int    `gorm:"column:tiktok;type:tinyint(1);default:0;NOT NULL" json:"tiktok"`                                        // 是否支持tiktok 0否 1是
	ForAll              int    `gorm:"column:for_all;type:tinyint(1);default:0;NOT NULL" json:"for_all"`                                      // 是否支持所有类型网站 0否 1是
	Banner              string `gorm:"column:banner;type:varchar(255)" json:"banner"`                                                         // 国旗
	IsDynamicThree      int    `gorm:"column:is_dynamic_three;type:tinyint(1);default:0;NOT NULL" json:"is_dynamic_three"`                    // 是否支持动态线路3 0否 1是
	Sneaker             int    `gorm:"column:sneaker;type:int(11)" json:"sneaker"`                                                            // 是否支持Sneaker 0否 1是
	StaticNum           int    `gorm:"column:static_num;type:int(11);default:0;NOT NULL" json:"static_num"`                                   // 静态住宅数量
	IproyalDcNum        int    `gorm:"column:iproyal_dc_num;type:int(11);default:0;NOT NULL" json:"iproyal_dc_num"`                           // 皇冠静态机房数量
	V6Num               int    `gorm:"column:v6_num;type:int(11);NOT NULL" json:"v6_num"`                                                     // 静态IPV6数量
	Proxy6DcNum         int    `gorm:"column:proxy6_dc_num;type:int(11);default:0;NOT NULL" json:"proxy6_dc_num"`                             // proxy6静态机房数量
	LineDcNum           int    `gorm:"column:line_dc_num;type:int(11);default:0;NOT NULL" json:"line_dc_num"`                                 // proxyline静态机房数量
	TransitStaticNum    int    `gorm:"column:transit_static_num;type:int(11);default:0;NOT NULL" json:"transit_static_num"`                   // 中转住宅数量
	TransitDcNum        int    `gorm:"column:transit_dc_num;type:int(11);default:0;NOT NULL" json:"transit_dc_num"`                           // 中转机房数量
	TransitV6Num        int    `gorm:"column:transit_v6_num;type:int(11);default:0;NOT NULL" json:"transit_v6_num"`                           // 中转V6数量
	TransitDc           int    `gorm:"column:transit_dc;type:tinyint(1);default:0;NOT NULL" json:"transit_dc"`                                // 是否开启中转机房 0否 1是
	TransitStatic       int    `gorm:"column:transit_static;type:tinyint(1);default:0;NOT NULL" json:"transit_static"`                        // 是否开启中转住宅 0否 1是
	TransitV6           int    `gorm:"column:transit_v6;type:tinyint(1);default:0;NOT NULL" json:"transit_v6"`                                // 是否开启中转V6 0否 1是
	SelfDc              int    `gorm:"column:self_dc;type:int(11);default:0;NOT NULL" json:"self_dc"`                                         // 自建机房数量
	SelfStatic          int    `gorm:"column:self_static;type:int(11);default:0;NOT NULL" json:"self_static"`                                 // 自建住宅数量
	KyDc                int    `gorm:"column:ky_dc;type:int(11)" json:"ky_dc"`                                                                // kookey机房id
	KyIsp               int    `gorm:"column:ky_isp;type:int(11)" json:"ky_isp"`                                                              // kookey 住宅id
}

func (m *VsIPCountryInfo) TableName() string {
	return VsIPCountryInfoTableName
}
