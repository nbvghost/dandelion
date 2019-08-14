package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;
import org.springframework.web.servlet.ModelAndView;

import javax.persistence.*;
import java.io.Serializable;
import java.text.SimpleDateFormat;
import java.util.Date;

/**
 * Created by SIX4 on 2014-10-18 .
 */
@Entity
@Table
public class Shop implements Serializable {
    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;
    private String business_name;
    private String branch_name;//分支
    private String province;
    private String city;
    private String district;
    private String address;
    private String telephone;
    private String categories;//门店的类型
    private Integer offset_type = 1;
    private String longitude;
    private String latitude;
    private String photo_list;
    private String special;
    private String open_time;
    private Integer avg_price;
    private String introduction;
    private String recommend;
    //private String poiId;
    //private String qrcode;
    private Boolean showLottery = false;
    private Boolean showYuyue = false;
    private Boolean showSeckill = false;
    //private String showCard;
    //@Column(columnDefinition = "varchar(255) default 'express:0'",nullable = false)
    private String authority;
    @Column(columnDefinition = "int default 0")
    private Integer state = 0;
    private Long expire;
    @Column(updatable = false)
    private Long date;
    private Integer vip;
    //private String topLevel;

    public Boolean getShowYuyue() {
        return showYuyue;
    }

    public void setShowYuyue(Boolean showYuyue) {
        this.showYuyue = showYuyue;
    }


    public Boolean getShowSeckill() {
        return showSeckill;
    }

    public void setShowSeckill(Boolean showSeckill) {
        this.showSeckill = showSeckill;
    }


    public String getBranch_name() {
        return branch_name;
    }

    public void setBranch_name(String branch_name) {
        this.branch_name = branch_name;
    }

    public String getProvince() {
        return province;
    }

    public void setProvince(String province) {
        this.province = province;
    }

    public String getCity() {
        return city;
    }

    public void setCity(String city) {
        this.city = city;
    }

    public String getDistrict() {
        return district;
    }

    public void setDistrict(String district) {
        this.district = district;
    }

    public String getCategories() {
        return categories;
    }

    public void setCategories(String categories) {
        this.categories = categories;
    }

    public Integer getOffset_type() {
        return offset_type;
    }

    public void setOffset_type(Integer offset_type) {
        this.offset_type = offset_type;
    }

    public String getLongitude() {
        return longitude;
    }

    public void setLongitude(String longitude) {
        this.longitude = longitude;
    }

    public String getLatitude() {
        return latitude;
    }

    public void setLatitude(String latitude) {
        this.latitude = latitude;
    }

    public String getPhoto_list() {
        return photo_list;
    }

    public void setPhoto_list(String photo_list) {
        this.photo_list = photo_list;
    }

    public String getSpecial() {
        return special;
    }

    public void setSpecial(String special) {
        this.special = special;
    }

    public String getOpen_time() {
        return open_time;
    }

    public void setOpen_time(String open_time) {
        this.open_time = open_time;
    }

    public Integer getAvg_price() {
        return avg_price;
    }

    public void setAvg_price(Integer avg_price) {
        this.avg_price = avg_price;
    }

    public String getIntroduction() {
        return introduction;
    }

    public void setIntroduction(String introduction) {
        this.introduction = introduction;
    }

    public String getRecommend() {
        return recommend;
    }

    public void setRecommend(String recommend) {
        this.recommend = recommend;
    }

    /*public String getPoiId() {
        return poiId;
    }

    public void setPoiId(String poiId) {
        this.poiId = poiId;
    }*/


    public Boolean getShowLottery() {
        return showLottery;
    }

    public void setShowLottery(Boolean showLottery) {
        this.showLottery = showLottery;
    }

    public Long getDate() {
        return date;
    }

    public void setDate(Long date) {
        this.date = date;
    }

    public String getBusiness_name() {
        return business_name;
    }

    public void setBusiness_name(String name) {
        this.business_name = name;
    }


    public String getTelephone() {
        return telephone;
    }

    public void setTelephone(String tel) {
        this.telephone = tel;
    }

    public String getAddress() {
        return address;
    }

    public void setAddress(String address) {
        this.address = address;
    }


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }


    public Long getExpire() {
        return expire;
    }

    public void setExpire(Long expire) {
        this.expire = expire;
    }

    public Integer getVip() {
        return vip;
    }

    public void setVip(Integer vip) {
        this.vip = vip;
    }

    /*public String getShowCard() {
        return showCard;
    }

    public void setShowCard(String showCard) {
        this.showCard = showCard;
    }*/

    public Integer getState() {
        return state;
    }

    public void setState(Integer state) {
        this.state = state;
    }

    /*public String getQrcode() {
        return qrcode;
    }

    public void setQrcode(String qrcode) {
        this.qrcode = qrcode;
    }*/

    /*public String getTopLevel() {
        return topLevel;
    }

    public void setTopLevel(String topLevel) {
        this.topLevel = topLevel;
    }*/

    public String getAuthority() {
        return authority;
    }

    public void setAuthority(String authority) {
        this.authority = authority;
    }
    public static class Authority{
        private Boolean express=false;

        public Boolean getExpress() {
            return express;
        }

        public void setExpress(Boolean express) {
            this.express = express;
        }
    }
    public static class State {
        public static final int Default = 0;
        public static final int Added = 1;
        public static final int Passed = 3;
    }

    public static class VIP {

        //0~1000
        public static final int leveA = 1000;
        //1001~2000
        public static final int leveB = 2000;
        //20001~3000
        public static final int leveC = 3000;
        public static final int not = -1;

        public static int PRICE_LEVE_A =9000;
        public static int PRICE_LEVE_B =18000;
        public static int PRICE_LEVE_C =36000;
        public static int EXPERIENCE = 12;//体验天数
        public static int REGISTER_REWARD = 50;//注册奖励，分
        public static int MAX_WITHDRAW = 1000;//最大提现金额，分

        public static int VIPDay(int leve) {
            switch (leve) {
                case leveA:
                    return 90;
                case leveB:
                    return 180;
                case leveC:
                    return 365;
                case not:
                    return EXPERIENCE;//体验12天时间
                default:
                    return 0;
            }
        }

        public static int isVip(int v) {
            if (v >= 0 && v <= 1000) {
                return leveA;
            } else if (v >= 1001 && v <= 2000) {
                return leveB;
            } else if (v >= 2001 && v <= 3000) {
                return leveC;
            } else {
                return -1;
            }
        }

        public static int VIPPrice(int leve) {
            /*switch (leve){
                case leveA:
                    return 9900;
                case leveB:
                    return 19900;
                case leveC:
                    return 33300;
                default:
                    return Integer.MAX_VALUE;
            }*/
            switch (leve) {
                case leveA:
                    return PRICE_LEVE_A;
                case leveB:
                    return PRICE_LEVE_B;
                case leveC:
                    return PRICE_LEVE_C;
                default:
                    return Integer.MAX_VALUE;

            }
        }

        public static String VIPName(int leve) {
            switch (leve) {
                case leveA:
                    return "黄金(三月)";
                case leveB:
                    return "白金(半年)";
                case leveC:
                    return "钻石(一年)";
                default:
                    return "超级会员(永久)";

            }
        }
        public static void Expire(Shop shop,Subscriber user, ModelAndView modelAndView){

            String msg = "";
            String url = "";
            String label = "";
            SimpleDateFormat dateFormater = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");

            ///account/order/${action}/${shopID}/item


            if (shop.getExpire() == 0 && shop.getVip() != Shop.VIP.not) {
                url = "/account/pay/platform_pay?shopID="+shop.getId()+"&action="+Orders.Action.create_upgrade()+"&openID="+user.getOpenID();
                msg = "您的账户即将到期(" + dateFormater.format(new Date(shop.getExpire())) + "到期)，请及时续费。";
                label = "续费";
            } else if (shop.getExpire() == 0 && shop.getVip() == Shop.VIP.not) {
                url = "/account/pay/platform_pay?shopID="+shop.getId()+"&action="+Orders.Action.create_register()+"&openID="+user.getOpenID();
                msg = "您的免费账户即将到期(" + dateFormater.format(new Date(shop.getExpire())) + "到期)，请及时升级。";
                label = "升级";
            } else if (shop.getExpire() != 0 && shop.getVip() == Shop.VIP.not) {
                url = "/account/pay/platform_pay?shopID="+shop.getId()+"&action="+Orders.Action.create_register()+"&openID="+user.getOpenID();
                msg = "您正在使用我们的免费账户，为了保证服务质量，请升级我们的会员账户。";
                label = "升级";
            } else {
                url = "/account/pay/platform_pay?shopID="+shop.getId()+"&action="+Orders.Action.create_upgrade()+"&openID="+user.getOpenID();
                msg = "您的账户即将到期(" + dateFormater.format(new Date(shop.getExpire())) + "到期)，请及时续费。";
                label = "续费";
            }
            //pay/platform_pay
            modelAndView.addObject("msg", msg);
            modelAndView.addObject("url", url);
            modelAndView.addObject("label", label);
        }
    }
}
