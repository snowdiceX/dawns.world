package dawns.twilight.dao.model;

import java.io.Serializable;
import java.util.Date;
import lombok.Data;

/**
 *
 * This class was generated by MyBatis Generator.
 * This class corresponds to the database table gaming
 *
 * @mbg.generated do_not_delete_during_merge Thu Feb 28 20:26:08 CST 2019
 */
@Data
public class Gaming implements Serializable {
    /**
     * gaming id
     */
    private Integer id;

    /**
     * game id
     */
    private Integer gameId;

    /**
     * chaincode id
     */
    private String gameContractId;

    /**
     * 
     */
    private String title;

    /**
     * 
     */
    private String description;

    /**
     * sequence of game　in progress
     */
    private String gamingSeq;

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
     * gaming
     */
    private static final long serialVersionUID = 1L;

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table gaming
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
        sb.append(", gameId=").append(gameId);
        sb.append(", gameContractId=").append(gameContractId);
        sb.append(", title=").append(title);
        sb.append(", description=").append(description);
        sb.append(", gamingSeq=").append(gamingSeq);
        sb.append(", status=").append(status);
        sb.append(", createtime=").append(createtime);
        sb.append(", remarks=").append(remarks);
        sb.append("]");
        return sb.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table gaming
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
        Gaming other = (Gaming) that;
        return (this.getId() == null ? other.getId() == null : this.getId().equals(other.getId()))
            && (this.getGameId() == null ? other.getGameId() == null : this.getGameId().equals(other.getGameId()))
            && (this.getGameContractId() == null ? other.getGameContractId() == null : this.getGameContractId().equals(other.getGameContractId()))
            && (this.getTitle() == null ? other.getTitle() == null : this.getTitle().equals(other.getTitle()))
            && (this.getDescription() == null ? other.getDescription() == null : this.getDescription().equals(other.getDescription()))
            && (this.getGamingSeq() == null ? other.getGamingSeq() == null : this.getGamingSeq().equals(other.getGamingSeq()))
            && (this.getStatus() == null ? other.getStatus() == null : this.getStatus().equals(other.getStatus()))
            && (this.getCreatetime() == null ? other.getCreatetime() == null : this.getCreatetime().equals(other.getCreatetime()))
            && (this.getRemarks() == null ? other.getRemarks() == null : this.getRemarks().equals(other.getRemarks()));
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table gaming
     *
     * @mbg.generated Thu Feb 28 20:26:08 CST 2019
     */
    @Override
    public int hashCode() {
        final int prime = 31;
        int result = 1;
        result = prime * result + ((getId() == null) ? 0 : getId().hashCode());
        result = prime * result + ((getGameId() == null) ? 0 : getGameId().hashCode());
        result = prime * result + ((getGameContractId() == null) ? 0 : getGameContractId().hashCode());
        result = prime * result + ((getTitle() == null) ? 0 : getTitle().hashCode());
        result = prime * result + ((getDescription() == null) ? 0 : getDescription().hashCode());
        result = prime * result + ((getGamingSeq() == null) ? 0 : getGamingSeq().hashCode());
        result = prime * result + ((getStatus() == null) ? 0 : getStatus().hashCode());
        result = prime * result + ((getCreatetime() == null) ? 0 : getCreatetime().hashCode());
        result = prime * result + ((getRemarks() == null) ? 0 : getRemarks().hashCode());
        return result;
    }
}