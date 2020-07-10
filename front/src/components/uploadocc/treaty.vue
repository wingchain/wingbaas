<template>
  <div class="create_contain">
    <div class="lince_title">合约管理</div>
    <!-- 添加用户对话框  :visible.sync 控制显示隐藏-->
    <!-- <div class="newCluster">
      <el-button type="primary" @click="addClose = true">+ 新建合约</el-button>
    </div>-->
    <div class="add_all">
      <div class="add_liance">
        <div class="add_box" @click="addClose = true">
          <div class="add_line">+</div>
          <div class="add_title">新建合约</div>
        </div>
      </div>
      <div class="linceList" v-for="(item, index) in allinceList" :key="index">
        <!-- {{allinceList}} -->
        <div class="add_box">
          <div class="topShow">
            <div class="brother">
              <div class="add_line heName">{{item.CCName}}</div>
              <div class="add_line lianName">{{item.BlockChainName}}</div>
            </div>
            <div class="top" @click="uploadLeave(item.CCName, item.CCVersion, item.BlockChainName)">升级</div>
          </div>
          <div class="add_title"></div>
          <div class="count_wrap">
            <div class="zuzhi">
              <span class="instru">版本</span>
              <div class="zuzhiCount">{{item.CCVersion}}</div>
            </div>
            <div class="jiedian">
              <span class="instru">更新时间</span>
              <div class="zuzhiCount">{{item.UpdateTime}}</div>
            </div>
            <div class="kuaigao">
              <span class="instru">创建时间</span>
              <div class="zuzhiCount">{{item.CreateTime}}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 新建合约 -->
    <el-dialog title="新建合约" :visible.sync="addClose" width="50%" class="add_lince">
      <!-- 内容主体区 -->
      <!-- model数据绑定，验证规则，引用对象 -->
      <el-form :model="addForm" :rules="addFormRules" ref="addFormRef" label-width="100px">
        <el-form-item label="合约名称" prop="ChainCodeId">
          <el-input v-model="addForm.ChainCodeId" placeholder="请输入2~20个英文，数字"></el-input>
        </el-form-item>
        <el-form-item label="联盟链" prop="BlockChainId">
          <el-select v-model="addForm.BlockChainId" placeholder="请选择联盟链">
            <el-option
              v-for="(item, index) in cluserData"
              :key="index"
              :label="item.BlockChainName"
              :value="item.BlockChainId"
            ></el-option>
            <!-- <el-option label="区域二" value="beijing"></el-option> -->
          </el-select>
        </el-form-item>
        <el-form-item prop="File" label="合约内容" ref="crtRef">
          <el-upload
            ref="crtRefREmove"
            class="upload-demo"
            action="http://106.75.51.138:9001/api/v1/uploadkeyfile"
            multiple
            :limit="1"
            accept=".go"
            :before-upload="beforeAvatarUpload"
            :on-success="succCrt"
          >
            <div class="upcrt">+ 上传合约</div>
          </el-upload>
        </el-form-item>
        <el-form-item label="版本号" prop="ChainCodeVersion">
          <el-input v-model="addForm.ChainCodeVersion" placeholder="请输入2~20个英文，数字"></el-input>
        </el-form-item>
      </el-form>
      <!-- 底部 -->
      <span slot="footer" class="dialog-footer">
        <el-button @click="addClose = false">取消</el-button>
        <el-button type="primary" @click="subCluster">提交</el-button>
      </span>
    </el-dialog>
    <!-- 合约升级-------------------------------------------------------------------------------- -->
    <el-dialog title="合约升级" :visible.sync="upgradeClose" width="50%" class="add_lince">
      <!-- 内容主体区 -->
      <!-- model数据绑定，验证规则，引用对象 -->
      <!-- :rules="addFormUpdataRules" -->
      <el-form
        :model="upgradeForm"
        ref="addFormLoadRef"
        label-width="100px"
        :rules="addForLoadmRef"
      >
        <el-form-item label="合约名称">
          <el-input v-model="upgradeForm.ChainCodeId" disabled></el-input>
        </el-form-item>
        <el-form-item label="联盟链">
          <el-input v-model="BlockName" disabled></el-input>
        </el-form-item>
        <el-form-item prop="File" label="合约内容" >
          <el-upload
          ref="crtTopRef"
            class="upload-demo"
            action="http://106.75.51.138:9001/api/v1/uploadkeyfile"
            multiple
            :limit="1"
            accept=".go"
            :before-upload="beforeUpdataUpload"
            :on-success="succTop"
          >
            <div class="upcrt">+ 上传合约</div>
          </el-upload>
        </el-form-item>
        <el-form-item label="版本号" prop="ChainCodeVersion">
          <el-input v-model="upgradeForm.ChainCodeVersion" placeholder="请输入2~20个数字"></el-input>
        </el-form-item>
      </el-form>
      <!-- 底部 -->
      <span slot="footer" class="dialog-footer">
        <el-button @click="upgradeClose = false">取消</el-button>
        <el-button type="primary" @click="subupLoad">提交</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
export default {
  data () {
    // 校验合约名称
    var validateName = (rule, value, callback) => {
      const mailReg = /^[a-zA-Z][a-zA-Z0-9-]{2,20}$/
      if (!value) {
        return callback(new Error('请输入合约名称'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('格式不正确，请输入2~20英文或数字的组合'))
        }
      }, 100)
    }
    // 校验版本号名称
    var validateVersion = (rule, value, callback) => {
      const mailReg = /^([1-9]\d|[1-9])(.([1-9]\d|\d)){2}$/
      if (!value) {
        return callback(new Error('请输入版本号'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('格式不正确，请重新输入'))
        }
      }, 100)
    }
    // 校验版本号名称
    var validateVersionLoad = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (this.upgradeForm.ChainCodeVersion === this.addForm.ChainCodeVersion) {
        callback(new Error('版本号重复'))
      }
      //  else {
      //   callback()
      // }
      const mailReg = /^([1-9]\d|[1-9])(.([1-9]\d|\d)){2}$/
      // if (value === '') {
      //   return callback(new Error('请输入版本号'))
      // } else if (value === this.addForm.ChainCodeVersion) {
      //   callback(new Error('版本号重复!'))
      // }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('格式不正确，请重新输入'))
        }
      }, 100)
    }
    return {
      addClose: false, // 新建联盟对话框状态
      upgradeClose: false, // 加入联盟对话框状态
      allinceList: [], // 合约列表数据
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
      formData: {},
      countList: {},
      lianList: [],
      cluData: [],
      AllianId: '',
      orgDo: '',
      BlockName: '',
      lian: [],
      isAdd: false,
      upgradeForm: {
        ChainCodeId: '',
        BlockChainId: '',
        ChainCodeVersion: '',
        File: ''
      },
      // 升级合约
      orgupData: {
        BlockChainId: '', //
        OrgName: '', //
        ChannelId: '',
        ChainCodeId: '', //
        ChainCodeVersion: '',
        ChaincodeSeq: '', // wu
        EndorsePolicy: '', // wu
        InitArgs: ['init', 'a1', '300', 'b1', '300']
      },
      // 新建合约的表单数据
      addForm: {
        BlockChainId: '',
        ChainCodeId: '',
        ChainCodeVersion: '',
        File: ''
      },
      // 部署合约
      yccData: {
        BlockChainId: '', // 联盟链id
        OrgName: '',
        ChannelId: 'mychannel',
        // ChannelId: 'baaschannel',
        ChainCodeId: '', // 合约名
        ChainCodeVersion: '', // 版本号
        EndorsePolicy: '', // 先为空
        InitArgs: ['init', 'a', '200', 'b', '200']
      },
      orgSingle: {
        BlockChainId: '',
        OrgName: '',
        ChannelId: 'mychannel',
        ChainCodeVersion: '',
        ChainCodeId: ''
      },
      // 合约数据
      queryForm: {
        BlockChainId: '',
        OrgName: '',
        ChannelId: 'mychannel'
      },
      // AllianceId: 'QlVwTJjMA3r90TMA',
      AllianceId: window.localStorage.getItem('AllianceId'),
      // 上传合约表单的验证规则
      addFormRules: {
        ChainCodeId: [
          {
            required: true,
            trigger: 'blur',
            validator: validateName
          }
        ],
        BlockChainId: [
          {
            required: true,
            trigger: 'blur',
            message: '请选择联盟链'
          }
        ],
        ChainCodeVersion: [
          {
            required: true,
            trigger: 'blur',
            validator: validateVersion,
            message: '请输入版本号'
          }
        ],
        File: [
          {
            required: true,
            trigger: 'blur',
            message: '请上传合约'
          }
        ]
      },
      // 升级合约校验
      addForLoadmRef: {
        ChainCodeVersion: [
          {
            required: true,
            trigger: 'blur',
            validator: validateVersionLoad
          }
        ],
        File: [
          {
            required: true,
            trigger: 'blur',
            message: '请上传合约'
          }
        ]
      }
    }
  },
  created () {
    this.getClur()
    // if (this.allinceList) {
    this.getjoinList()
    // }
  },
  methods: {
    // 升级
    async uploadLeave (name, version, cluName) {
      this.upgradeClose = true
      this.upgradeForm.ChainCodeId = name
      this.addForm.ChainCodeVersion = version
      this.BlockName = cluName
    },
    // 加入的合约列表
    async getjoinList () {
      this.allinceList = []
      const { data: resClu } = await this.$http.get(`${this.AllianceId}/alliancechains`)
      this.cluData = resClu.data
      this.queryForm.OrgName = window.sessionStorage.getItem('OrgName')
      for (var h = 0; h < this.cluData.length; h++) {
        this.queryForm.BlockChainId = this.cluData[h].BlockChainId
        // this.lian.push(resName.data)
        // console.log('连名称组合', this.lian)
        const { data: res } = await this.$http.post('queryinstatialcc', this.queryForm)
        console.log('请求来的数据', res.data)
        const { data: resName } = await this.$http.get(`${this.queryForm.BlockChainId}/blockchain`)
        console.log('链名称', resName.data.BlockChainName)
        /*
         1: 直接拉取合约列表：根据blockChaniId 第一条链 + 第二条链的数据
         2：新增合约拉取列表: 追加之前，先将之前的置空
         */
        this.allinceList = [...this.allinceList, ...res.data]
      }
      // 如果当前合约的某个值和拉取链名字的值一样，就加name
      console.log('此刻的数据', this.allinceList)
    },
    // 新建合约上传证书
    beforeAvatarUpload (file) {
      this.addForm.File = file
      this.formData = new FormData()
      this.formData.append('file', file)
      // this.$refs.crtRef.clearValidate()
    },
    succCrt () {
      console.log(this.$refs.crtRef)
      this.$refs.crtRef.clearValidate()
    },
    // 升级合约上传证书
    beforeUpdataUpload (file) {
      this.upgradeForm.File = file
      this.formData = new FormData()
      this.formData.append('file', file)
      // this.$refs.crtTopRef.clearValidate()
    },
    succTop () {
      this.$refs.crtTopRef.clearValidate()
    },
    // succ () {

    // },
    async getClur () {
      const { data: res } = await this.$http.get(`${this.AllianceId}/alliancechains`)
      if (res.code !== 0) return this.$message.error('拉取链列表失败')
      //   this.$message.success('拉取链列表成功')
      this.cluserData = res.data
      console.log('链数据', this.cluserData)
    },
    // 上传合约
    async subCluster () {
      this.$refs.addFormRef.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        console.log('预校验', this.addForm)
        this.formData.append('BlockChainId', this.addForm.BlockChainId)
        this.formData.append('ChainCodeId', this.addForm.ChainCodeId)
        this.formData.append('ChainCodeVersion', this.addForm.ChainCodeVersion)
        console.log('合约信息', this.formData.get('file'))
        const { data: res } = await this.$http.post('uploadcc', this.formData)
        console.log('上传合约列表数据', res)
        if (res.code !== 0) return this.$message.error('上传合约失败')
        const { data: cfg } = await this.$http.get(`${this.addForm.BlockChainId}/cfg`)
        if (cfg.code !== 0) return this.$message.error('拉取节点组织失败')
        // this.$message.success('拉取节点组织成功')
        this.countList = cfg.data
        console.log('节点数据', this.countList)
        // this.$message.success('上传合约成功')
        this.yccData.ChainCodeVersion = this.addForm.ChainCodeVersion
        this.yccData.ChainCodeId = this.addForm.ChainCodeId
        this.yccData.BlockChainId = this.addForm.BlockChainId
        this.yccData.OrgName = this.countList.DeployNetCfg.PeerOrgs[0].Name
        // this.yccData.BlockChainId = this.addForm.BlockChainId
        // window.sessionStorage.setItem('BlockChainId', this.addForm.BlockChainId)
        // window.sessionStorage.setItem('ChannelId', this.yccData.ChannelId)
        window.sessionStorage.setItem('OrgName', this.countList.DeployNetCfg.PeerOrgs[0].Name)
        const { data: ycc } = await this.$http.post('orgdeploycc', this.yccData)
        console.log('部署合约列表数据', ycc)
        if (ycc.code !== 0) return this.$message.error('部署合约失败')
        this.$message.success('部署合约成功')
        this.isAdd = true
        this.getjoinList()
        this.addClose = false
        this.$refs.addFormRef.resetFields()
        if (this.addForm.File) {
          console.log(this.addForm.File)
          console.log(this.$refs.crtRef)
          this.$refs.crtRefREmove.clearFiles()
        }
        // 隐藏添加用户的对话框
      })
    },
    // 升级合约
    async subupLoad () {
      this.$refs.addFormLoadRef.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        console.log('预校验', this.queryForm)
        // this.getCount()
        this.formData.append('ChainCodeId', this.upgradeForm.ChainCodeId) // 合约名称
        this.formData.append('BlockChainId', this.queryForm.BlockChainId) // 联盟链
        this.formData.append('ChainCodeVersion', this.upgradeForm.ChainCodeVersion) // 版本
        console.log('合约信息', this.formData.get('file'))
        const { data: res } = await this.$http.post('uploadcc', this.formData)
        console.log('合约升级上传合约代码数据', res)
        if (res.code !== 0) return this.$message.error('上传合约失败')
        const { data: cfg } = await this.$http.get(`${this.queryForm.BlockChainId}/cfg`)
        if (cfg.code !== 0) return this.$message.error('拉取节点组织失败')
        this.$message.success('拉取节点组织成功')
        this.countList = cfg.data
        console.log('节点数据', this.countList)
        // 调singleorgdeploycc--------------------------------------------------------------
        for (var j = 0; j < this.countList.DeployNetCfg.PeerOrgs.length; j++) {
          console.log('节点数据', this.countList.DeployNetCfg.PeerOrgs)
          this.orgSingle.BlockChainId = this.queryForm.BlockChainId
          this.orgSingle.OrgName = this.countList.DeployNetCfg.PeerOrgs[j].Name
          this.orgSingle.ChainCodeVersion = this.upgradeForm.ChainCodeVersion
          this.orgSingle.ChainCodeId = this.upgradeForm.ChainCodeId
          const { data: chunk } = await this.$http.post('singleorgdeploycc', this.orgSingle)
          console.log('合约成功第二步', chunk)
          if (chunk.code !== 0) return this.$message.error('合约失败')
        }
        // 合约升级-------------------------------------------------------------------------
        this.orgupData.BlockChainId = this.queryForm.BlockChainId // 联盟链
        this.orgupData.OrgName = this.queryForm.OrgName // 组织
        this.orgupData.ChainCodeVersion = this.upgradeForm.ChainCodeVersion // 合约版本
        this.orgupData.ChainCodeId = this.upgradeForm.ChainCodeId // 合约名称
        this.orgupData.ChannelId = this.queryForm.ChannelId // ChannelId
        const { data: upLeave } = await this.$http.post('orgupgradecc', this.orgupData)
        console.log('升级合约列表数据', upLeave)
        if (upLeave.code !== 0) return this.$message.error('升级合约失败')
        this.$message.success('升级合约成功')
        this.isAdd = true
        this.getjoinList()
        this.upgradeForm.ChainCodeVersion = ''
        if (this.upgradeForm.File) {
          this.$refs.crtTopRef.clearFiles()
        }
        // 隐藏添加用户的对话框
        this.upgradeClose = false
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
  .add_all {
    display: flex;
    .add_liance {
      margin-top: 27px;
      width: 300px;
      height: 133px;
      background-color: #f6f8fb;
      border-radius: 12px;
      border: solid 1px #c8c9cc;
      .add_box {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 100%;
        font-size: 18px;
        color: #27afb5;
        cursor: pointer;
        .add_line {
          margin-right: 5px;
        }
      }
    }
    .linceList {
      margin-left: 22px;
      margin-top: 27px;
      width: 450px;
      height: 133px;
      background-color: #f6f8fb;
      border-radius: 12px;
      border: solid 1px #c8c9cc;
      .del_llince {
        color: #c8c9cc;
        font-size: 20px;
        margin-left: 115px;
        margin-top: 0px;
        cursor: pointer;
      }
      .add_box {
        padding: 14px;
        height: 100%;
        font-size: 18px;
        color: #1f232d;
        cursor: pointer;
        .topShow {
          display: flex;
          justify-content: space-between;
          .brother {
            .heName {
              font-size: 20px;
            }
            .lianName {
              background: #dfebff;
              color: #107bff;
              display: flex;
              justify-content: center;
              align-items: center;
              margin-top: 5px;
              font-size: 12px;
            }
          }
          .top {
            color: #27afb5;
          }
          .add_title {
            font-size: 14px;
          }
        }
        .count_wrap {
          display: flex;
          margin-top: 12px;
          .jiedian,
          .kuaigao {
            text-align: center;
            width: 43%;
          }
          .zuzhi {
            width: 13%;
          }
          .zuzhiCount {
            font-size: 14px;
            color: #1f232d;
            margin-top: 2px;
          }
          .instru {
            color: #717580;
            font-size: 14px;
          }
        }
      }
    }
  }
  .newCluster {
    font-size: 9px;
    margin-top: 10px;
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
}
</style>
