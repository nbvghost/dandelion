package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by sixf on 2/19/2016.
 */
@Entity
@Table
public class Systems implements Serializable {
    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;
    @Column(length = 10240)
    private String value;
    private String label;
    private String des;

    public String getLabel() {
        return label;
    }

    public void setLabel(String label) {
        this.label = label;
    }
    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDes() {
        return des;
    }

    public void setDes(String des) {
        this.des = des;
    }

    public static class SystemLabel {

        public static final String WX_MENU_LABEL = "WX_MENU_LABEL";
        public static final String SECKILL_DATE_LABEL = "SECKILL_DATE_LABEL";

        public static final String PRICE_LEVE_A = "PRICE_LEVE_A";
        public static final String PRICE_LEVE_B = "PRICE_LEVE_B";
        public static final String PRICE_LEVE_C = "PRICE_LEVE_C";
        public static final String EXPERIENCE = "EXPERIENCE";

        public static final String REGISTER_REWARD = "REGISTER_REWARD";
        public static final String MAX_WITHDRAW = "MAX_WITHDRAW";
    }
}
