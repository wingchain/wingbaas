<template>
  <div class="create_contain">
    <div class="lince_title">联盟机构列表</div>
    <!-- 添加用户对话框  :visible.sync 控制显示隐藏-->
    <div class="newCluster">
      <el-button type="primary" @click="joinClose = true">+ 邀请机构</el-button>
    </div>
    <el-table :data="surprice" border style="width:100%" stripe>
      <el-table-column prop="Mail" label="机构名称" width="500">
        <template slot-scope="scope">
          {{scope.row.Mail}}
          <span
            v-for="(itemSub, indexSub) in scope.row.Alliances"
            :key="indexSub"
          >
            <!-- {{scope.row.Alliances}}
            && itemSub.Mail === Mail-->
            <span
              v-if="itemSub.Id === addForm.AllianceId && itemSub.Creator === scope.row.Mail"
              class="Honst"
            >盟主</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="加入时间" width="500" prop>
        <template slot-scope="scope">
          <!-- {{scope.row.Mail}} -->
          <span v-for="(itemSub, indexSub) in scope.row.Alliances" :key="indexSub">
            <!-- {{scope.row.Alliances}}
            && itemSub.Mail === Mail-->
            <span v-if="itemSub.Id === addForm.AllianceId">{{itemSub.JoinTime}}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="level" label="操作" width="281" class="caozuo">
        <template slot-scope="scope">
          <!-- {{scope.row}} -->
          <!-- <span @click="del_user(scope.row.Mail, scope.row.Alliances[0].Id)" style="color:#ff1743">删除</span> -->
          <span @click="del_user(addForm.AllianceId, scope.row.Mail)" class="del_usr">删除</span>
        </template>
      </el-table-column>
    </el-table>
    <!-- 邀请结构 -->
    <el-dialog title="邀请机构" :visible.sync="joinClose" width="50%" class="add_linceVane">
      <el-form
        :model="joinForm"
        :rules="addVanceFormRules"
        ref="addVanceFormRef"
        label-width="162px"
      >
        <el-form-item prop="Mail" label="受邀机构电子邮箱">
          <el-input v-model="joinForm.Mail" style="margin-left:4px" placeholder="请输入你收到的联盟邀请码"></el-input>
        </el-form-item>
      </el-form>
      <!-- 底部 -->
      <span slot="footer" class="dialog-footer">
        <el-button @click="joinLince(joinForm.Alliance.Id)" type="primary">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
export default {
  data () {
    // 校验邮箱
    var checkEmail = (rule, value, callback) => {
      const mailReg = /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/
      if (!value) {
        return callback(new Error('请输入邮箱'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('邮箱格式不正确，请重新输入'))
        }
      }, 100)
    }
    return {
      addClose: false, // 新建联盟对话框状态
      joinClose: false, // 加入联盟对话框状态
      allinceList: [], // 联盟列表数据
      Mail: '', //   取来的邮箱
      crt: '', // 证书
      key: '', // 私钥
      clusId: '', // 传入集群的id
      cluserData: [], // 集群列表数据
      DataAll: '',
      cluId: '',
      clData: [],
      newData: [],
      surprice: [],
      AllianceId: '',
      // 邀请用户表单数据
      joinForm: {
        Mail: '',
        Alliance: {
          Name: '',
          Describe: '',
          Creator: '',
          Id: window.localStorage.getItem('AllianceId')
        }
      },
      addVanceFormRules: {
        Mail: [
          {
            required: true,
            validator: checkEmail,
            trigger: 'blur'
          }
        ]
      },
      // 添加用户的表单数据
      addForm: {
        AllianceId: window.localStorage.getItem('AllianceId')
      }
    }
  },
  created () {
    this.Mail = window.sessionStorage.getItem('Mail')
    this.AllianceId = window.localStorage.getItem('AllianceId')
    console.log('关键字', this.AllianceId)
    // 拉取加入的联盟用户列表
    this.getJoinList()
  },
  methods: {
    async getJoinList () {
      const { data: res } = await this.$http.get(`${this.AllianceId}/users`)
      this.surprice = res.data
      console.log('获取加入的用户', this.surprice)
    },
    // 删除用户
    async del_user (id, Mail) {
      const delForm = {
        Mail: Mail,
        AllianceId: id
      }
      const isConfird = await this.$confirm(`此操作将永久删除『${Mail}』, 是否继续`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      console.log(isConfird)
      if (isConfird !== 'confirm') {
        return this.$message.info('已经取消删除')
      }
      const { data: res } = await this.$http.post('deleteallianceuser', delForm)
      if (res.code !== 0) return this.$message.error('删除用户失败')
      this.$message.success('删除用户成功')
      this.getJoinList()
    },
    // 邀请加入联盟
    async joinLince (id) {
      this.$refs.addVanceFormRef.validate(async valid => {
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        this.addForm.Creator = this.Mail
        console.log('预校验', this.addForm)
        const { data: res } = await this.$http.get(`${id}/alliance`)
        console.log('拉取联盟数据', res.data)
        if (res.code !== 0) return this.$message.error('加入联盟失败')
        this.joinForm.Alliance = res.data
        console.log(this.joinForm)
        const { data: resAdd } = await this.$http.post('useraddalliance', this.joinForm)
        if (resAdd.code !== 0) return this.$message.error('邀请用户失败')
        this.$message.success('邀请用户成功')
        this.joinClose = false
        console.log('加入联盟数据', resAdd)
        this.getJoinList()
        this.joinForm.Mail = ''
      })
    }
  }
}
</script>
<style lang="less" scoped>
.create_contain {
  position: relative;
  width: 100%;
  .lince_title {
    font-size: 24px;
    color: #1f232d;
  }
  .newCluster {
    font-size: 9px;
    margin-top: 27px;
    /deep/ .el-button--primary {
      box-shadow: 3px 4px 16px 0px rgba(39, 71, 152, 0.32);
      border-radius: 4px;
    }
  }
  .el-form {
    .upcrt {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 95px;
      height: 28px;
      background-color: #ebf2ff;
      border-radius: 14px;
      border: solid 1px #2e82ff;
      cursor: pointer;
    }
  }
  .box-card {
    margin-top: 1rem;
  }
  .my_create {
    font-size: 14px;
    color: #1f232d;
    display: flex;
    align-items: center;
    margin-top: 30px;
    .line {
      width: 4px;
      height: 10px;
      background-color: #27afb5;
      margin-right: 7px;
    }
  }
  /deep/ .el-table--border td {
    border-right: none;
  }
  /deep/ .el-table td,
  .el-table th.is-leaf {
    border-bottom: none;
  }
  /deep/ .el-table th.is-leaf {
    border-bottom: none;
  }
  /deep/ .el-table--border th {
    border-bottom: none;
    border-right: none;
    background: #f6f8fb;
  }
  /deep/ .el-table--border {
    border-radius: 12px;
    // border: solid 1px #c8c9cc;
  }
  /deep/ .el-table--striped .el-table__body tr.el-table__row--striped td {
    background: #f6f8fb;
  }
  /deep/ .el-table__body {
    width: 100% !important;
  }
  /deep/ .el-table__header {
    width: 1200% !important;
  }
  .del_usr {
    color: #ff1743;
    cursor: pointer;
  }
  .Honst {
    background: #2e82ff;
    padding: 5px 8px;
    color: #333;
    border-radius: 5px;
  }
  // /deep/ .el-table th:nth-child(3) {
  //   text-align: right;
  // }
  // /deep/ .el-table__row .el-table .cell {
  //   text-align: right;
  // }
}
</style>
