package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by sixf on 3/12/2016.
 */
@Entity
@Table
public class Orders implements Serializable {
    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;
    private String payNo;
    @Column(updatable = false)
    private String shopID;
    @Column(updatable = false)
    private Integer amount;//分计算
    private String name;
    private String des;
    @Column(updatable = false)
    private Integer type=-1;
    @Column(updatable = false)
    private Integer action;
    private Long payDate;
    @Column(updatable = false)
    private Long createDate;
    private String status;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }


    public Integer getAmount() {
        return amount;
    }

    public void setAmount(Integer amount) {
        this.amount = amount;
    }

    public Long getCreateDate() {
        return createDate;
    }

    public void setCreateDate(Long createDate) {
        this.createDate = createDate;
    }

    public Long getPayDate() {
        return payDate;
    }

    public void setPayDate(Long payDate) {
        this.payDate = payDate;
    }

    public String getDes() {
        return des;
    }

    public void setDes(String des) {
        this.des = des;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Integer getType() {
        return type;
    }

    public void setType(Integer type) {
        this.type = type;
    }


    public String getPayNo() {
        return payNo;
    }

    public void setPayNo(String payNo) {
        this.payNo = payNo;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }


    public Integer getAction() {
        return action;
    }

    public void setAction(Integer action) {
        this.action = action;
    }

    public String getShopID() {
        return shopID;
    }

    public void setShopID(String shopID) {
        this.shopID = shopID;
    }

    public static class Action {
        public static final int register = 100;//注册交费
        public static final int register_reward = 200;//注册奖励
        public static final int renew = 1000;//续费
        public static final int transfers = 10000;//提现
        public static final int express = 20000;//快递收款
        public static final int mall = 30000;//商城
        public static final int mall_one = 40000;//一元购

        //0~100
        public static int create_register() {
            return (int) (Math.random() * 100);
        }

        //101~1000
        public static int create_upgrade() {
            int s = (int) (Math.random() * 900);
            if (s == 0) {
                s = 1;
            }
            return 100 + s;
        }

        //1001~10000
        public static int create_transfers() {
            int s = (int) (Math.random() * 9000);
            if (s == 0) {
                s = 1;
            }
            return 1000 + s;
        }

        public static int isAction(int v) {
            if (v >= 0 && v <= 100) {
                return register;
            } else if (v >= 101 && v <= 1000) {
                return renew;
            } else if (v >= 1001 && v <= 10000) {
                return transfers;
            }else if (v == 20000) {
                return express;
            } else {
                return -1;
            }
        }
    }

    public static class Status {
        public static final String paying = "paying";
        public static final String payed = "payed";
    }
}
