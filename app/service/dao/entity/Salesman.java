package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by sixf on 2016/7/28.
 */
@Entity
@Table
public class Salesman implements Serializable{
    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;
    private String pass;
    private Boolean validate=false;
    private String shopID;
    private Long date;
    /**

     * cascade = CascadeType.MERGE--级联更新 cascade=CascadeType.PERSIST--级联持久
     * cascade=CascadeType.REMOVE--级联删除 cascade=CascadeType.REFRESH-- 级联刷新
     * 在业务逻辑中可能对象进行修改，但是读取出来并不是最新的数据。 如果需要最新的数据，这时就得需要级联刷新
     * fetch = FetchType.LAZY--开启延迟加载。 fetch = FetchType.EAGER--即时加载
     * 在数据中，这个字段是否为空 optional=false,这个选项不可以空

     */
    @ManyToOne(cascade = CascadeType.REFRESH,fetch = FetchType.EAGER)
    @JoinColumn(name = "userID")
    private Subscriber user;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public Boolean getValidate() {
        return validate;
    }

    public void setValidate(Boolean validate) {
        this.validate = validate;
    }

    public Long getDate() {
        return date;
    }

    public void setDate(Long date) {
        this.date = date;
    }

    public String getShopID() {
        return shopID;
    }

    public void setShopID(String shopID) {
        this.shopID = shopID;
    }

    public Subscriber getUser() {
        return user;
    }

    public void setUser(Subscriber user) {
        this.user = user;
    }

    public String getPass() {
        return pass;
    }

    public void setPass(String pass) {
        this.pass = pass;
    }
}
