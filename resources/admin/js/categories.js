/**
 *Createdbysixfon3/15/2016.
 */

function Categories() {
    this.Items = {};
}
Categories.prototype.add = function (id, iArray) {
    this.Items[id] = iArray;
}
Categories.prototype.Exists = function (id) {
    if (typeof(this.Items[id]) == "undefined"){
        return false;
    }
    return true;
}
function CategoriesCode() {
    this.Items = {};
}
CategoriesCode.prototype.add = function (id, iArray) {
    this.Items[id] = iArray;
}
CategoriesCode.prototype.Exists = function (id) {
    if (typeof(this.Items[id]) == "undefined"){
        return false;
    }
    return true;
}

var categories = new Categories();
var categoriesCode = new CategoriesCode();


categories.add("0", ["美食", "基础设施", "医疗保健", "生活服务", "休闲娱乐", "购物", "运动健身", "汽车", "酒店宾馆", "旅游景点", "文化场馆", "教育学校", "银行金融", "地名地址", "房产小区", "丽人", "结婚", "亲子", "公司企业", "机构团体", "其它"]);

categories.add("0_0", ["江浙菜", "粤菜", "川菜", "湘菜", "东北菜", "徽菜", "闽菜", "鲁菜", "台湾菜", "西北菜", "东南亚菜", "西餐", "日韩菜", "火锅", "清真菜", "小吃快餐", "海鲜", "烧烤", "自助餐", "面包甜点", "茶餐厅", "咖啡厅", "其它美食"]);
categories.add("0_1", ["交通设施", "公共设施", "道路附属", "其它基础设施"]);
categories.add("0_2", ["专科医院", "综合医院", "诊所", "急救中心", "药房药店", "疾病预防", "其它医疗保健"]);
categories.add("0_3", ["家政", "宠物服务", "旅行社", "摄影冲印", "洗衣店", "票务代售", "邮局速递", "通讯服务", "彩票", "报刊亭", "自来水营业厅", "电力营业厅", "教练", "生活服务场所", "信息咨询中心", "招聘求职", "中介机构", "事务所", "丧葬", "废品收购站", "福利院养老院", "测字风水", "家装", "其它生活服务"]);
categories.add("0_4", ["洗浴推拿足疗", "KTV", "酒吧", "咖啡厅", "茶馆", "电影院", "棋牌游戏", "夜总会", "剧场音乐厅", "度假疗养", "户外活动", "网吧", "迪厅", "演出票务", "其它娱乐休闲"]);
categories.add("0_5", ["综合商场", "便利店", "超市", "花鸟鱼虫", "家具家居建材", "体育户外", "服饰鞋包", "图书音像", "眼镜店", "母婴儿童", "珠宝饰品", "化妆品", "食品烟酒", "数码家电", "农贸市场", "小商品市场", "旧货市场", "商业步行街", "礼品", "摄影器材", "钟表店", "拍卖典当行", "古玩字画", "自行车专卖", "文化用品", "药店", "品牌折扣店", "其它购物"]);
categories.add("0_6", ["健身中心", "游泳馆", "瑜伽", "羽毛球馆", "乒乓球馆", "篮球场", "足球场", "壁球场", "马场", "高尔夫场", "保龄球馆", "溜冰", "跆拳道", "海滨浴场", "网球场", "橄榄球", "台球馆", "滑雪", "舞蹈", "攀岩馆", "射箭馆", "综合体育场馆", "其它运动健身"]);
categories.add("0_7", ["加油站", "停车场", "4S店", "汽车维修", "驾校", "汽车租赁", "汽车配件销售", "汽车保险", "摩托车", "汽车养护", "洗车场", "汽车俱乐部", "汽车救援", "二手车交易市场", "车辆管理机构", "其它汽车"]);
categories.add("0_8", ["星级酒店", "经济型酒店", "公寓式酒店", "度假村", "农家院", "青年旅社", "酒店宾馆", "旅馆招待所", "其它酒店宾馆"]);
categories.add("0_9", ["公园", "其它旅游景点", "风景名胜", "植物园", "动物园", "水族馆", "城市广场", "世界遗产", "国家级景点", "省级景点", "纪念馆", "寺庙道观", "教堂", "海滩"]);
categories.add("0_10", ["博物馆", "图书馆", "美术馆", "展览馆", "科技馆", "天文馆", "档案馆", "文化宫", "会展中心", "其它文化场馆"]);
categories.add("0_11", ["小学", "幼儿园", "其它教育学校", "培训", "大学", "中学", "职业技术学校", "成人教育"]);
categories.add("0_12", ["银行", "自动提款机", "保险公司", "证券公司", "财务公司", "其它银行金融"]);
categories.add("0_13", ["交通地名", "地名地址信息", "道路名", "自然地名", "行政地名", "门牌信息", "其它地名地址"]);
categories.add("0_14", ["住宅区", "产业园区", "商务楼宇", "它房产小区"]);
categories.add("0_15", ["美发", "美容", "SPA", "瘦身纤体", "美甲", "写真", "其它"]);
categories.add("0_16", ["婚纱摄影", "婚宴", "婚戒首饰", "婚纱礼服", "婚庆公司", "彩妆造型", "司仪主持", "婚礼跟拍", "婚车租赁", "婚礼小商品", "婚房装修", "其它"]);
categories.add("0_17", ["亲子摄影", "亲子游乐", "亲子购物", "孕产护理"]);
categories.add("0_18", ["农林牧渔基地", "企业/工厂", "其它公司企业"]);
categories.add("0_19", ["公检法机构", "外国机构", "工商税务机构", "政府机关", "民主党派", "社会团体", "传媒机构", "文艺团体", "科研机构", "其它机构团体"]);
categories.add("0_20", ["其它"]);


categories.add("0_0_0", ["上海菜", "淮扬菜", "浙江菜", "南京菜", "苏帮菜", "杭帮菜", "宁波菜", "无锡菜", "舟山菜", "衢州菜", "绍兴菜", "温州菜", "苏北土菜"]);
categoriesCode.add("0_0_0", ["101010", "101011", "101012", "101013", "101014", "101015", "101016", "101017", "101018", "101019", "10101a", "10101b", "10101c"]);

categories.add("0_0_1", ["潮汕菜", "茶餐厅", "客家菜", "湛江菜"]);
categoriesCode.add("0_0_1", ["101110", "101111", "101112", "101113"]);

categories.add("0_0_2", ["自贡盐帮菜", "江湖菜", "酸菜鱼", "香锅", "川味小吃"]);
categoriesCode.add("0_0_2", ["101210", "101211", "101212", "101213", "101214"]);

categoriesCode.add("0_0_3", ["101300"]);
categoriesCode.add("0_0_4", ["101400"]);
categoriesCode.add("0_0_5", ["101500"]);
categoriesCode.add("0_0_6", ["101600"]);
categoriesCode.add("0_0_7", ["101700"]);
categoriesCode.add("0_0_8", ["101800"]);
categoriesCode.add("0_0_9", ["101900"]);

categories.add("0_0_10", ["泰国菜", "越南菜", "印度菜", "星马菜", "其它东南亚菜"]);
categoriesCode.add("0_0_10", ["101a10", "101a11", "101a12", "101a13", "101aff"]);

categories.add("0_0_11", ["法国菜", "意大利菜", "俄罗斯菜", "牛排", "比萨", "巴西菜", "中东菜", "西式正餐", "西式简餐", "西班牙菜", "无国界料理", "美国菜", "葡国菜", "地中海菜", "拉美烧烤", "英国菜", "德国菜", "墨西哥菜", "其它西餐"]);
categoriesCode.add("0_0_11", ["101b10", "101b11", "101b12", "101b13", "101b14", "101b15", "101b16", "101b17", "101b18", "101b19", "101b1a", "101b1b", "101b1c", "101b1d", "101b1e", "101b1f", "101b20", "101b21", "101bff"]);

categories.add("0_0_12", ["日本菜", "韩国菜", "其它日韩菜"]);
categoriesCode.add("0_0_12", ["101c10", "101c11", "101cff"]);


categoriesCode.add("0_0_13", ["101d00"]);
categoriesCode.add("0_0_14", ["101e00"]);
categoriesCode.add("0_0_15", ["101f00"]);
categoriesCode.add("0_0_16", ["102000"]);
categoriesCode.add("0_0_17", ["102100"]);
categoriesCode.add("0_0_18", ["102200"]);
categoriesCode.add("0_0_19", ["102300"]);
categoriesCode.add("0_0_20", ["102400"]);
categoriesCode.add("0_0_20", ["102500"]);

categories.add("0_0_22", ["北京家常菜", "官府菜", "云贵菜", "湖北菜", "山西菜", "豫菜", "天津菜", "新疆菜", "农家菜", "创意菜", "素菜", "烤鸭", "江西菜", "蒙古菜", "广西菜", "冀菜", "陕西菜", "青海菜", "西藏菜", "酒楼", "家常菜", "私家菜", "民族菜", "冷饮店", "其它中餐厅"]);
categoriesCode.add("0_0_22", ["10ff10", "10ff11", "10ff12", "10ff13", "10ff14", "10ff15", "10ff16", "10ff17", "10ff18", "10ff19", "10ff1a", "10ff1b", "10ff1c", "10ff1d", "10ff1e", "10ff1f", "10ff20", "10ff21", "10ff22", "10ff23", "10ff24", "10ff25", "10ff26", "10ff27", "10ffff"]);


categories.add("0_1_0", ["交通服务相关", "公交车站", "地铁站", "港口码头", "火车站", "轻轨站", "过境口岸", "长途汽车站", "飞机场", "公交线路", "地铁线路", "其它交通设施"]);
categoriesCode.add("0_1_0", ["111010", "111011", "111012", "111013", "111014", "111015", "111016", "111017", "111018", "111019", "11101a", "1110ff"]);

categories.add("0_1_1", ["公共厕所", "公用电话", "紧急避难场所", "其它公共设施"]);
categoriesCode.add("0_1_1", ["111110", "111111", "111112", "1111ff"]);

categories.add("0_1_2", ["收费站", "服务区", "其它道路附属"]);
categoriesCode.add("0_1_2", ["111210", "111211", "1112ff"]);

categories.add("0_1_3", ["其它基础设施"]);
categoriesCode.add("0_1_3", ["11ff00"]);


categories.add("0_2_0", ["齿科", "整形", "眼科", "耳鼻喉", "胸科", "骨科", "肿瘤", "脑科", "妇产科", "儿科", "传染病医院", "精神病医院", "其它专科医院"]);
categoriesCode.add("0_2_0", ["121010", "121011", "121012", "121013", "121014", "121015", "121016", "121017", "121018", "121019", "12101a", "12101b", "1210ff"]);

categoriesCode.add("0_2_1", ["121100"]);
categoriesCode.add("0_2_2", ["121200"]);
categoriesCode.add("0_2_3", ["121300"]);
categoriesCode.add("0_2_4", ["121400"]);
categoriesCode.add("0_2_5", ["121500"]);
categoriesCode.add("0_2_6", ["12ff00"]);


categories.add("0_3_0", ["月嫂保姆", "保洁钟点工", "开锁", "送水", "家电维修", "管道疏通打孔", "搬家", "其它家政"]);
categoriesCode.add("0_3_0", ["131010", "131011", "131012", "131013", "131014", "131015", "131016", "1310ff"]);

categories.add("0_3_1", ["宠物商店", "宠物市场", "宠物医院", "其它宠物服务"]);
categoriesCode.add("0_3_1", ["131110", "131111", "131112", "1311ff"]);

categoriesCode.add("0_3_2", ["131200"]);
categoriesCode.add("0_3_3", ["131300"]);
categoriesCode.add("0_3_4", ["131400"]);

categories.add("0_3_5", ["飞机票代售", "火车票代售", "汽车票代售", "公交及IC卡", "景点售票", "其它票务代售"]);
categoriesCode.add("0_3_5", ["131510", "131511", "131512", "131513", "131514", "1315ff"]);

categories.add("0_3_6", ["邮局", "速递"]);
categoriesCode.add("0_3_6", ["131610", "131611"]);

categories.add("0_3_7", ["中国电信营业厅", "中国网通营业厅", "中国移动营业厅", "中国联通营业厅", "中国铁通营业厅", "其它通讯服务"]);
categoriesCode.add("0_3_7", ["131710", "131711", "131712", "131713", "131714", "1317ff"]);

categories.add("0_3_8", ["彩票彩券销售点", "马会投注站"]);

categoriesCode.add("0_3_8", ["131810", "131811"]);
categoriesCode.add("0_3_9", ["131900"]);
categoriesCode.add("0_3_10", ["131a00"]);
categoriesCode.add("0_3_11", ["131b00"]);
categoriesCode.add("0_3_12", ["131c00"]);
categoriesCode.add("0_3_13", ["131d00"]);
categoriesCode.add("0_3_14", ["131e00"]);
categoriesCode.add("0_3_15", ["131f00"]);
categoriesCode.add("0_3_16", ["132000"]);
categoriesCode.add("0_3_17", ["132100"]);
categoriesCode.add("0_3_18", ["132200"]);
categoriesCode.add("0_3_19", ["132300"]);
categoriesCode.add("0_3_20", ["132400"]);
categoriesCode.add("0_3_21", ["132500"]);
categoriesCode.add("0_3_22", ["132600"]);
categoriesCode.add("0_3_23", ["13ff00"]);


categoriesCode.add("0_4_0", ["141000"]);
categoriesCode.add("0_4_1", ["141100"]);
categoriesCode.add("0_4_2", ["141200"]);
categoriesCode.add("0_4_3", ["141300"]);
categoriesCode.add("0_4_4", ["141400"]);
categoriesCode.add("0_4_5", ["141500"]);
categoriesCode.add("0_4_6", ["141600"]);
categoriesCode.add("0_4_7", ["141700"]);
categoriesCode.add("0_4_8", ["141800"]);
categoriesCode.add("0_4_9", ["141900"]);
categoriesCode.add("0_4_10", ["141a00"]);
categoriesCode.add("0_4_11", ["141b00"]);
categoriesCode.add("0_4_12", ["141c00"]);
categoriesCode.add("0_4_13", ["141d00"]);
categoriesCode.add("0_4_14", ["14ff00"]);


categoriesCode.add("0_5_0", ["151000"]);
categoriesCode.add("0_5_1", ["151100"]);
categoriesCode.add("0_5_2", ["151200"]);
categoriesCode.add("0_5_3", ["151300"]);
categoriesCode.add("0_5_4", ["151400"]);
categoriesCode.add("0_5_5", ["151500"]);
categoriesCode.add("0_5_6", ["151600"]);
categoriesCode.add("0_5_7", ["151700"]);
categoriesCode.add("0_5_8", ["151800"]);
categoriesCode.add("0_5_9", ["151900"]);
categoriesCode.add("0_5_10", ["151a00"]);
categoriesCode.add("0_5_11", ["151b00"]);
categoriesCode.add("0_5_12", ["151c00"]);
categoriesCode.add("0_5_13", ["151d00"]);
categoriesCode.add("0_5_14", ["151e00"]);
categoriesCode.add("0_5_15", ["151f00"]);
categoriesCode.add("0_5_16", ["152000"]);
categoriesCode.add("0_5_17", ["152100"]);
categoriesCode.add("0_5_18", ["152200"]);
categoriesCode.add("0_5_19", ["152300"]);
categoriesCode.add("0_5_20", ["152400"]);
categoriesCode.add("0_5_21", ["152500"]);
categoriesCode.add("0_5_22", ["152600"]);
categoriesCode.add("0_5_23", ["152700"]);
categoriesCode.add("0_5_24", ["152800"]);
categoriesCode.add("0_5_25", ["152900"]);
categoriesCode.add("0_5_26", ["152a00"]);
categoriesCode.add("0_5_27", ["15ff00"]);


categoriesCode.add("0_6_0", ["161000"]);
categoriesCode.add("0_6_1", ["161100"]);
categoriesCode.add("0_6_2", ["161200"]);
categoriesCode.add("0_6_3", ["161300"]);
categoriesCode.add("0_6_4", ["161400"]);
categoriesCode.add("0_6_5", ["161500"]);
categoriesCode.add("0_6_6", ["161600"]);
categoriesCode.add("0_6_7", ["161700"]);
categoriesCode.add("0_6_8", ["161800"]);
categoriesCode.add("0_6_9", ["161900"]);
categoriesCode.add("0_6_10", ["161a00"]);
categoriesCode.add("0_6_11", ["161b00"]);
categoriesCode.add("0_6_12", ["161c00"]);
categoriesCode.add("0_6_13", ["161d00"]);
categoriesCode.add("0_6_14", ["161e00"]);
categoriesCode.add("0_6_15", ["161f00"]);
categoriesCode.add("0_6_16", ["162000"]);
categoriesCode.add("0_6_17", ["162100"]);
categoriesCode.add("0_6_18", ["162200"]);
categoriesCode.add("0_6_19", ["162300"]);
categoriesCode.add("0_6_20", ["162400"]);
categoriesCode.add("0_6_21", ["162500"]);
categoriesCode.add("0_6_22", ["16ff00"]);


categoriesCode.add("0_7_0", ["171000"]);
categoriesCode.add("0_7_1", ["171100"]);
categoriesCode.add("0_7_2", ["171200"]);
categoriesCode.add("0_7_3", ["171300"]);
categoriesCode.add("0_7_4", ["171400"]);
categoriesCode.add("0_7_5", ["171500"]);
categoriesCode.add("0_7_6", ["171600"]);
categoriesCode.add("0_7_7", ["171700"]);
categoriesCode.add("0_7_8", ["171800"]);
categoriesCode.add("0_7_9", ["171900"]);
categoriesCode.add("0_7_10", ["171a00"]);
categoriesCode.add("0_7_11", ["171b00"]);
categoriesCode.add("0_7_12", ["171c00"]);
categoriesCode.add("0_7_13", ["171d00"]);
categoriesCode.add("0_7_14", ["171e00"]);
categoriesCode.add("0_7_15", ["17ff00"]);


categoriesCode.add("0_8_0", ["181000"]);
categoriesCode.add("0_8_1", ["181100"]);
categoriesCode.add("0_8_2", ["181200"]);
categoriesCode.add("0_8_3", ["181300"]);
categoriesCode.add("0_8_4", ["181400"]);
categoriesCode.add("0_8_5", ["181500"]);
categoriesCode.add("0_8_6", ["181600"]);
categoriesCode.add("0_8_7", ["181700"]);
categoriesCode.add("0_8_8", ["18ff00"]);


categoriesCode.add("0_9_0", ["191000"]);
categoriesCode.add("0_9_1", ["19ff00"]);
categoriesCode.add("0_9_2", ["191100"]);
categoriesCode.add("0_9_3", ["191200"]);
categoriesCode.add("0_9_4", ["191300"]);
categoriesCode.add("0_9_5", ["191400"]);
categoriesCode.add("0_9_6", ["191500"]);
categoriesCode.add("0_9_7", ["191600"]);
categoriesCode.add("0_9_8", ["191700"]);
categoriesCode.add("0_9_9", ["191800"]);
categoriesCode.add("0_9_10", ["191900"]);
categoriesCode.add("0_9_11", ["191a00"]);
categoriesCode.add("0_9_12", ["191b00"]);
categoriesCode.add("0_9_13", ["191c00"]);


categoriesCode.add("0_10_0", ["1a1000"]);
categoriesCode.add("0_10_1", ["1a1100"]);
categoriesCode.add("0_10_2", ["1a1200"]);
categoriesCode.add("0_10_3", ["1a1300"]);
categoriesCode.add("0_10_4", ["1a1400"]);
categoriesCode.add("0_10_5", ["1a1500"]);
categoriesCode.add("0_10_6", ["1a1600"]);
categoriesCode.add("0_10_7", ["1a1700"]);
categoriesCode.add("0_10_8", ["1a1800"]);
categoriesCode.add("0_10_9", ["1aff00"]);


categoriesCode.add("0_11_0", ["1b1000"]);
categoriesCode.add("0_11_1", ["1b1100"]);
categoriesCode.add("0_11_2", ["1bff00"]);
categoriesCode.add("0_11_3", ["1b1200"]);
categoriesCode.add("0_11_4", ["1b1300"]);
categoriesCode.add("0_11_5", ["1b1400"]);
categoriesCode.add("0_11_6", ["1b1500"]);
categoriesCode.add("0_11_7", ["1b1600"]);


categoriesCode.add("0_12_0", ["1c1000"]);
categoriesCode.add("0_12_1", ["1c1100"]);
categoriesCode.add("0_12_2", ["1c1200"]);
categoriesCode.add("0_12_3", ["1c1300"]);
categoriesCode.add("0_12_4", ["1c1400"]);
categoriesCode.add("0_12_5", ["1cff00"]);


categoriesCode.add("0_13_0", ["1d1000"]);
categoriesCode.add("0_13_1", ["1d1100"]);
categoriesCode.add("0_13_2", ["1d1200"]);
categoriesCode.add("0_13_3", ["1d1300"]);
categoriesCode.add("0_13_4", ["1d1400"]);
categoriesCode.add("0_13_5", ["1d1500"]);
categoriesCode.add("0_13_6", ["1dff00"]);


categoriesCode.add("0_14_0", ["1e1000"]);
categoriesCode.add("0_14_1", ["1e1100"]);
categoriesCode.add("0_14_2", ["1e1200"]);
categoriesCode.add("0_14_3", ["1eff00"]);


categoriesCode.add("0_15_0", ["1f1000"]);
categoriesCode.add("0_15_1", ["1f1000"]);
categoriesCode.add("0_15_2", ["1f1200"]);
categoriesCode.add("0_15_3", ["1f1300"]);
categoriesCode.add("0_15_4", ["1f1400"]);
categoriesCode.add("0_15_5", ["1f1500"]);
categoriesCode.add("0_15_6", ["1fff00"]);


categoriesCode.add("0_16_0", ["201000"]);
categoriesCode.add("0_16_1", ["201100"]);
categoriesCode.add("0_16_2", ["201200"]);
categoriesCode.add("0_16_3", ["201300"]);
categoriesCode.add("0_16_4", ["201400"]);
categoriesCode.add("0_16_5", ["201500"]);
categoriesCode.add("0_16_6", ["201600"]);
categoriesCode.add("0_16_7", ["201700"]);
categoriesCode.add("0_16_8", ["201800"]);
categoriesCode.add("0_16_9", ["201900"]);
categoriesCode.add("0_16_10", ["201a00"]);
categoriesCode.add("0_16_11", ["20ff00"]);


categoriesCode.add("0_17_0", ["211000"]);
categoriesCode.add("0_17_1", ["211100"]);
categoriesCode.add("0_17_2", ["211200"]);
categoriesCode.add("0_17_3", ["211300"]);


categoriesCode.add("0_18_0", ["221000"]);
categoriesCode.add("0_18_1", ["221100"]);
categoriesCode.add("0_18_2", ["22ff00"]);


categoriesCode.add("0_19_0", ["231000"]);
categoriesCode.add("0_19_1", ["231100"]);
categoriesCode.add("0_19_2", ["231200"]);
categoriesCode.add("0_19_3", ["231300"]);
categoriesCode.add("0_19_4", ["231400"]);
categoriesCode.add("0_19_5", ["231500"]);
categoriesCode.add("0_19_6", ["231600"]);
categoriesCode.add("0_19_7", ["231700"]);
categoriesCode.add("0_19_8", ["231800"]);
categoriesCode.add("0_19_9", ["23ff00"]);


categoriesCode.add("0_20_0", ["ff0000"]);