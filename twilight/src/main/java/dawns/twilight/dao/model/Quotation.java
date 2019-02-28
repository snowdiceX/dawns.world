package dawns.twilight.dao.model;

import java.io.Serializable;
import java.util.Date;
import lombok.Data;

/**
 *
 * This class was generated by MyBatis Generator.
 * This class corresponds to the database table quotation
 *
 * @mbg.generated do_not_delete_during_merge Thu Feb 28 20:26:08 CST 2019
 */
@Data
public class Quotation implements Serializable {
    /**
     * token quotation id
     */
    private Integer id;

    /**
     * 
     */
    private String baseTokenName;

    /**
     * 
     */
    private String network;

    /**
     * 
     */
    private String tokenName;

    /**
     * 
     */
    private Integer price;

    /**
     * 
     */
    private Integer status;

    /**
     * 
     */
    private Date createtime;

    /**
     * 
     */
    private String remarks;

    /**
     * quotation
     */
    private static final long serialVersionUID = 1L;

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table quotation
     *
     * @mbg.generated Thu Feb 28 20:26:08 CST 2019
     */
    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append(getClass().getSimpleName());
        sb.append(" [");
        sb.append("Hash = ").append(hashCode());
        sb.append(", id=").append(id);
        sb.append(", baseTokenName=").append(baseTokenName);
        sb.append(", network=").append(network);
        sb.append(", tokenName=").append(tokenName);
        sb.append(", price=").append(price);
        sb.append(", status=").append(status);
        sb.append(", createtime=").append(createtime);
        sb.append(", remarks=").append(remarks);
        sb.append("]");
        return sb.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table quotation
     *
     * @mbg.generated Thu Feb 28 20:26:08 CST 2019
     */
    @Override
    public boolean equals(Object that) {
        if (this == that) {
            return true;
        }
        if (that == null) {
            return false;
        }
        if (getClass() != that.getClass()) {
            return false;
        }
        Quotation other = (Quotation) that;
        return (this.getId() == null ? other.getId() == null : this.getId().equals(other.getId()))
            && (this.getBaseTokenName() == null ? other.getBaseTokenName() == null : this.getBaseTokenName().equals(other.getBaseTokenName()))
            && (this.getNetwork() == null ? other.getNetwork() == null : this.getNetwork().equals(other.getNetwork()))
            && (this.getTokenName() == null ? other.getTokenName() == null : this.getTokenName().equals(other.getTokenName()))
            && (this.getPrice() == null ? other.getPrice() == null : this.getPrice().equals(other.getPrice()))
            && (this.getStatus() == null ? other.getStatus() == null : this.getStatus().equals(other.getStatus()))
            && (this.getCreatetime() == null ? other.getCreatetime() == null : this.getCreatetime().equals(other.getCreatetime()))
            && (this.getRemarks() == null ? other.getRemarks() == null : this.getRemarks().equals(other.getRemarks()));
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table quotation
     *
     * @mbg.generated Thu Feb 28 20:26:08 CST 2019
     */
    @Override
    public int hashCode() {
        final int prime = 31;
        int result = 1;
        result = prime * result + ((getId() == null) ? 0 : getId().hashCode());
        result = prime * result + ((getBaseTokenName() == null) ? 0 : getBaseTokenName().hashCode());
        result = prime * result + ((getNetwork() == null) ? 0 : getNetwork().hashCode());
        result = prime * result + ((getTokenName() == null) ? 0 : getTokenName().hashCode());
        result = prime * result + ((getPrice() == null) ? 0 : getPrice().hashCode());
        result = prime * result + ((getStatus() == null) ? 0 : getStatus().hashCode());
        result = prime * result + ((getCreatetime() == null) ? 0 : getCreatetime().hashCode());
        result = prime * result + ((getRemarks() == null) ? 0 : getRemarks().hashCode());
        return result;
    }
}