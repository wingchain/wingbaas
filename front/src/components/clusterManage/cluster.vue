<template>
  <div class="create_contain">
    <div class="lince_title">集群管理</div>
    <!-- 添加用户对话框  :visible.sync 控制显示隐藏-->
    <div class="newCluster">
      <el-button type="primary" @click="addClose = true">+ 新建集群</el-button>
    </div>
    <el-dialog title="新建集群" :visible.sync="addClose" width="50%" class="add_lince">
      <!-- 内容主体区 -->
      <!-- model数据绑定，验证规则，引用对象 -->
      <el-form :model="addForm" :rules="addFormRules" ref="addFormRef" label-width="100px">
        <el-form-item label="集群名称" prop="ClusterId">
          <el-input v-model="addForm.ClusterId" placeholder="请输入集群名称"></el-input>
        </el-form-item>
        <el-form-item prop="PublicIp" label="公网Ip">
          <el-input v-model="addForm.PublicIp"></el-input>
        </el-form-item>
        <el-form-item prop="HostDomain" label="集群域名">
          <el-input v-model="addForm.HostDomain"></el-input>
        </el-form-item>
        <el-form-item label="ApiUrl" prop="ApiUrl">
          <el-input v-model="addForm.ApiUrl"></el-input>
        </el-form-item>
        <el-form-item label="集群证书" prop="Cert">
          <el-upload
            ref="crtRef"
            class="upload-demo"
            action="http://106.75.51.138:9001/api/v1/uploadkeyfile"
            multiple
            :on-success="handleAvatarSuccess"
            :limit="1"
            accept=".crt"
          >
            <div class="upcrt">+ 上传证书</div>
          </el-upload>
        </el-form-item>
        <el-form-item label="集群私钥" prop="Key">
          <el-upload
            ref="keyRef"
            class="upload-demo"
            action="http://106.75.51.138:9001/api/v1/uploadkeyfile"
            multiple
            :limit="1"
            :on-success="handleSshSuccess"
            accept=".key"
          >
            <div class="upcrt">+ 上传私钥</div>
          </el-upload>
        </el-form-item>
      </el-form>
      <!-- 底部 -->
      <span slot="footer" class="dialog-footer">
        <el-button @click="addClose = false">取消</el-button>
        <el-button type="primary" @click="subCluster">提交</el-button>
      </span>
    </el-dialog>
    <div v-for="(item, index) in surprice" :key="index">
      <div class="my_create">
        <div class="line"></div>集群名称
      </div>
      <el-table :data="item" border style="width:100%" stripe>
        <el-table-column prop="metadata.name" label="主机名称" width="500" >
          <template slot-scope="scope">
<span v-html="scope.row.metadata.name" style="color: #1f232d"></span>
          </template>
        </el-table-column>
        <el-table-column label="主机IP" width="500">
          <template slot-scope="scope">
            <span v-html="regFormate(scope.row.metadata.name)" style="color: #1f232d">
              <!-- <span v-html="scope.row.metadata.name"> -->
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="状态" width="265">
          <!-- 进行中 -->
          <template slot-scope="scope">
            <span v-if="scope.row.metadata.name" style="color: #27afb5">进行中</span>
            <span v-else>中断</span>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="item.length > 4">查看更多</div>
    </div>
  </div>
</template>
<script>
export default {
  data () {
    // 校验集群名称
    var validateName = (rule, value, callback) => {
      const mailReg = /^[a-zA-Z][a-zA-Z0-9-]{2,20}$/
      if (!value) {
        return callback(new Error('请输入集群名称'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('格式不正确，请输入2~20英文或数字的组合'))
        }
      }, 100)
    }
    // 校验集群名称
    var validateHost = (rule, value, callback) => {
      const mailReg = /^[a-zA-Z][a-zA-Z0-9-]{2,20}$/
      if (!value) {
        return callback(new Error('请输入集群域名'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('格式不正确，请输入2~20英文或数字的组合'))
        }
      }, 100)
    }
    // 校验api url名称
    var validateUrl = (rule, value, callback) => {
      const mailReg = /(http|https):\/\/([\w.]+\/?)\S*/
      if (!value) {
        return callback(new Error('请输入API URL'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('格式不正确，请输入以http或https开头的URL'))
        }
      }, 100)
    }
    // 校验公网IP
    var validateIp = (rule, value, callback) => {
      const mailReg = /^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$/
      if (!value) {
        return callback(new Error('请输入公网IP'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('格式不正确,请输入IP地址'))
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
      // 添加用户的表单数据
      addForm: {
        AllianceId: window.localStorage.getItem('AllianceId'),
        ClusterId: '',
        // Creator: this.Mail
        ApiUrl: 'https://kubernetes:6443',
        HostDomain: 'kubernetes',
        PublicIp: '106.75.51.138',
        Cert: '', // 集群证书
        Key: '', // 集群私钥
        InterIp: '172.16.254.33'
      },
      // 表单的验证规则
      addFormRules: {
        ApiUrl: [
          {
            required: true,
            trigger: 'blur',
            validator: validateUrl
          }
        ],
        HostDomain: [
          {
            required: true,
            message: '请输入集群域名',
            trigger: 'blur',
            validator: validateHost
          }
        ],
        PublicIp: [
          {
            required: true,
            trigger: 'blur',
            validator: validateIp
          }
        ],
        ClusterId: [
          {
            required: true,
            trigger: 'blur',
            validator: validateName
          }
        ],
        Cert: [
          {
            required: true,
            message: '请上传证书',
            trigger: 'blur'
          }
        ],
        Key: [
          {
            required: true,
            message: '请上传私钥',
            trigger: 'blur'
          }
        ]
      }
    }
  },
  created () {
    this.Mail = window.sessionStorage.getItem('Mail')
    this.getClur()
    // this.getclusterList()
    if (this.clusId) {
      this.getclusterList()
    }
  },
  methods: {
    regFormate (value) {
      var reg = value.replace(/-/gi, '.')
      return reg
    },
    // beforeAvatarUpload (file) {
    //   this.addForm.Cert = file
    // },
    // beforeUpdataUpload (file) {
    //   this.addForm.Key = file
    // },
    handleAvatarSuccess (response, file, fileList) {
      console.log('文件信息', response)
      this.addForm.Cert = response.data
      this.$refs.crtRef.clearValidate()
    },
    handleSshSuccess (response, file, fileList) {
      console.log('文件信息', response)
      this.addForm.Key = response.data
      this.$refs.keyRef.clearValidate()
    },
    async getClur () {
      const { data: res } = await this.$http.get('clusters')
      console.log('得到集群数据id', res)
      this.cluserData = res.data
      console.log(this.cluserData)
      for (var i = 0; i < this.cluserData.length; i++) {
        this.clusId = this.cluserData[i].ClusterId
        console.log('？？', this.clusId) // 新建集群时候的集群名字
        this.getclusterList()
      }
    },
    // 获取集群列表
    async getclusterList () {
      const { data: res } = await this.$http.get(`${this.clusId}/hosts`)
      console.log('获取集群列表', res)
      this.clData = res.data.items
      console.log('item项', this.clData)
      this.surprice.push(this.clData)
      console.log('填充进来的所有数据', this.surprice)
    },
    // 新建集群
    async subCluster () {
      this.$refs.addFormRef.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        console.log('预校验', this.addForm)
        const { data: res } = await this.$http.post('addcluster', this.addForm)
        console.log('新建联盟创建集群列表数据', res)
        if (res.code !== 0) return this.$message.error('新建联盟创建集群列表失败')
        this.$message.success('新建联盟创建集群列表成功')
        this.addForm.ClusterId = ''
        if (this.addForm.Cert) {
          console.log('数据', this.addForm.Cert)
          this.$refs.crtRef.clearFiles()
        }
        if (this.addForm.Key) {
          // this.$refs.keyRef.clearFiles()
        }
        // 隐藏添加用户的对话框
        this.addClose = false
        this.getclusterList()
      })
    }
  }
}
</script>
<style lang="less" scoped>
.create_contain {
  position: relative;
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
    width: 1210% !important;
  }
  // /deep/ .el-table .cell {
  //   color: #1f232d;
  // }
}
</style>
