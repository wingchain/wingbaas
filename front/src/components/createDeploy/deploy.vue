<template>
  <div class="create_contain">
    <div class="lince_title">区块链管理</div>
    <!-- 添加用户对话框  :visible.sync 控制显示隐藏-->
    <div class="newCluster">
      <el-button type="primary" @click="addClose = true">+ 新建联盟链</el-button>
    </div>
    <div class="lian" v-for="(item, index) in replaceList" :key="index">
      <div class="lian_containt">
        <div class="lian_wrap">
          <div class="lian_title">
            <div class="title">{{item.BlockChainName}}</div>
            <div class="more" @click="removeLian(item.BlockChainName, item.ClusterId)">删除</div>
          </div>
          <div class="version">版本号 {{item.BlockChainType}} &nbsp; {{item.Version}}</div>
        </div>
        <div class="count_wrap">
          <div class="zuzhi">
            <span>组织数(个)</span>
            <div class="zuzhiCount">{{item.zuzhi}}</div>
          </div>
          <div class="jiedian">
            <span>节点数(个)</span>
            <div class="zuzhiCount">{{item.jieCount}}</div>
          </div>
          <div class="kuaigao">
            <span>块高</span>
            <div class="zuzhiCount">{{item.heightCount}}</div>
          </div>
          <div class="yewu">
            <span>业务总量</span>
            <div class="zuzhiCount">{{item.TxCount}}</div>
          </div>
        </div>
      </div>
    </div>
    <el-dialog title="新建区块链" :visible.sync="addClose" width="50%" class="add_lince">
      <!-- 内容主体区 -->
      <!-- model数据绑定，验证规则，引用对象 -->
      <!-- .DeployCfg.DeployNetCfg -->
      <el-form :model="addForm" :rules="addFormRules" ref="addFormRef" label-width="110px">
        <el-form-item prop="BlockChainName" label="联盟链名称">
          <el-input v-model="addForm.BlockChainName" mimlength="2" maxlength="20" show-word-limit></el-input>
        </el-form-item>
        <el-form-item label="集群" prop="DeployCfg.ClusterId">
          <el-select v-model="addForm.DeployCfg.ClusterId" placeholder="请选择集群">
            <el-option
              v-for="(item, index) in cluserData"
              :key="index"
              :label="item.ClusterId"
              :value="item.ClusterId"
            ></el-option>
            <!-- <el-option label="区域二" value="beijing"></el-option> -->
          </el-select>
        </el-form-item>
        <div class="addPlug" @click="gojiqun">+ 新建集群</div>
        <!-- 动态增加表数据  一行的增加，三处和表单的校验 -->
        <el-form-item label="组织与节点配置">
          <div
            v-for="(item, index) in addForm.DeployCfg.DeployNetCfg.PeerOrgs"
            :key="index"
            class="pointAll"
          >
            <el-form-item
              class="plug"
              :prop="'DeployCfg.DeployNetCfg.PeerOrgs.' + index + '.Name'"
              :rules="DeployCfg.DeployNetCfg.PeerOrgs"
            >
              <el-input placeholder="组织名称，建议格式为xxxOrg" v-model="item.Name"></el-input>
            </el-form-item>
            <el-form-item
              v-for="(itemSub, indexSub) in addForm.DeployCfg.DeployNetCfg.PeerOrgs[index].Specs"
              :prop="'DeployCfg.DeployNetCfg.PeerOrgs.' + index + '.Specs.'+ indexSub + '.Hostname'"
              :key="indexSub"
              :rules="[{ validator: (rule, value, callback) => {
                const mailReg = /^[a-z][a-z-\d]+$/
      if (!value) {
        return callback('节点名不能为空')
      }
        if (mailReg.test(value)) {
          callback()
        } else {
          callback('请输入正确的节点名称格式')
        }
              }}]"
            >
              <el-input class="plugPoint" placeholder="节点名称" v-model="itemSub.Hostname"></el-input>
            </el-form-item>
            <el-form-item>
              <i class="el-icon-delete" @click="deleteItem(index)"></i>
            </el-form-item>
          </div>
        </el-form-item>
        <div class="addPlug anothor" @click="removeDomain">+ 新增组织</div>
        <div class="info">
          组织名称需以字母开头，12位以内的数字和字母组成。只有网络初创者可以新建
          多个组织。
        </div>
      </el-form>
      <!-- 底部 -->
      <span slot="footer" class="dialog-footer">
        <el-button @click="addClose = false">取消</el-button>
        <el-button type="primary" @click="subCluster">提交</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
export default {
  data () {
    var validateName = (rule, value, callback) => {
      const mailReg = /^[a-z][a-z-\d]+$/
      if (!value) {
        return callback(new Error('区块链不能为空'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('请输入正确的区块链名称格式'))
        }
      }, 100)
    }
    return {
      addClose: false, // 新建联盟对话框状态
      joinClose: false, // 加入联盟对话框状态
      allinceList: [], // 联盟列表数据
      clusId: '', // 传入集群的id
      AllianId: '',
      zuzhi: 0,
      jiedian: '',
      jieCount: 0,
      heightCount: 0,
      TxCount: 0,
      countList: [],
      heightList: [],
      lianList: [],
      cluserData: [], // 集群列表数据
      DataAll: '',
      cluId: '',
      clData: [],
      replaceList: [],
      org: '',
      orgDo: '',
      isHei: true,
      // 快高业务总量数据
      heightData: {
        BlockChainId: '',
        ChannelId: 'mychannel',
        OrgName: ''
      },
      paramsChanal: {
        BlockChainId: '',
        OrgName: '',
        ChannelId: 'mychannel'
      },
      orgChannel: {
        BlockChainId: '',
        OrgName: '',
        ChannelId: 'mychannel'
      },
      // 添加用户的表单数据
      addForm: {
        // BlockChainName: 'test-chainnetwork11', // 联盟链名称
        BlockChainName: '', // 联盟链名称
        BlockChainType: 'fabric',
        // AllianceId: 'QlVwTJjMA3r90TMA',
        AllianceId: window.localStorage.getItem('AllianceId'),
        DeployCfg: {
          DeployNetCfg: {
            OrdererOrgs: [
              {
                Name: 'Orderer',
                Domain: 'orderer.baas.xyz',
                DeployNode: '172-16-254-130',
                Specs: [
                  {
                    Hostname: 'orderer0'
                  },
                  {
                    Hostname: 'orderer1'
                  }
                ]
              }
            ],
            PeerOrgs: [
              // 名称
              {
                Name: 'Org1', // 组织
                // Name: '', // 组织
                Domain: '', // 域名
                DeployNode: '172-16-254-33',
                Specs: [
                  {
                    // Hostname: 'peer0-org1' // 节点
                    Hostname: '' // 节点
                  }
                ],
                Users: {
                  Count: 1
                }
              },
              {
                Name: 'Org2',
                // Name: '',
                Domain: '',
                DeployNode: '172-16-254-130',
                Specs: [
                  {
                    // Hostname: 'peer0-org2'
                    Hostname: ''
                  }
                ],
                Users: {
                  Count: 1
                }
              }
            ],
            KafkaDeployNode: '172-16-254-130',
            ZookeeperDeployNode: '172-16-254-33',
            ToolsDeployNode: '172-16-254-33'
          },
          DeployType: 'KAFKA_FABRIC',
          Version: '1.3.0',
          CryptoType: 'ECDSA',
          ClusterId: '' // 集群名称
        }
      },
      // 新增表单的验证规则
      DeployCfg: {
        DeployNetCfg: {
          PeerOrgs: [
            { required: true, message: '请输入组织名称', trigger: 'blur' },
            {
              validator: (rule, value, callback) => {
                // console.log('AAA', rule)
                // console.log('value', value)
                const reg = /^[A-Z][a-z-\d]+$/
                if (reg.test(value)) {
                  callback()
                } else {
                  callback(new Error('组织名称不符合要求'))
                }
                // if (value === this.DeployCfg.DeployNetCfg.PeerOrgs.splice(value.index)) {
                //   callback(new Error('组织名称重复'))
                // }
              },
              trigger: 'change'
            }
          ]
        }
      },
      // 表单的验证规则
      addFormRules: {
        //   联盟链名称
        BlockChainName: [
          {
            required: true,
            message: '请输入2~20英文，数字',
            trigger: 'blur',
            validator: validateName
          }
        ],
        HostDomain: [
          {
            required: true,
            message: '请输入集群域名',
            trigger: 'blur'
          }
        ],
        PublicIp: [
          {
            required: true,
            message: '请输入公网IP',
            trigger: 'blur'
          }
        ],
        DeployCfg: {
          ClusterId: [
            {
              required: true,
              message: '请选择集群',
              trigger: 'change'
            }
          ]
        }
      }
    }
  },
  created () {
    this.Mail = window.sessionStorage.getItem('Mail')
    this.getClur()
    this.getLianList()
  },
  methods: {
    // 拉取链列表
    async getLianList () {
      const { data: res } = await this.$http.get(`${this.addForm.AllianceId}/alliancechains`)
      if (res.code !== 0) return this.$message.error('拉取链列表失败')
      // this.$message.success('拉取链列表成功')
      this.lianList = res.data
      console.log('链列表数据', this.lianList)
      if (this.lianList === null) {
        this.replaceList = []
      }
      this.getCount()
      this.getHeight()
    },
    gojiqun () {
      this.$router.push('/cluster')
    },
    async getHeight () {
      if (this.lianList) {
        for (var T = 0; T < this.lianList.length; T++) {
          this.heightData.BlockChainId = this.lianList[T].BlockChainId
          this.heightData.OrgName = window.sessionStorage.getItem('OrgName')
          const { data: resHeight } = await this.$http.post('blocktx', this.heightData)
          if (resHeight.code !== 0) return this.$message.eror('拉取节点高度，总量失败')
          // this.$message.success('拉取节点高度，总量成功')
          this.heightList = resHeight.data
          if (this.heightList) {
            this.heightCount = this.heightList.Height
            this.TxCount = this.heightList.TxCount
          }
          this.$set(this.lianList[T], 'heightCount', this.heightCount)
          this.$set(this.lianList[T], 'TxCount', this.TxCount)
          this.replaceList = this.lianList
          this.heightCount = 0
          this.TxCount = 0
          console.log('push成功没有', this.replaceList)
        }
      }
    },
    async getCount () {
      if (this.lianList) {
        for (var A = 0; A < this.lianList.length; A++) {
          this.AllianId = this.lianList[A].BlockChainId
          const { data: res } = await this.$http.get(`${this.AllianId}/cfg`)
          // this.$message.success('拉取节点组织成功')
          this.countList = res.data
          this.zuzhi = this.countList.DeployNetCfg.PeerOrgs.length // 组织个数
          for (var w = 0; w < this.countList.DeployNetCfg.PeerOrgs.length; w++) {
            this.jieCount = this.jieCount + this.countList.DeployNetCfg.PeerOrgs[w].Specs.length
          }
          this.$set(this.lianList[A], 'zuzhi', this.zuzhi)
          this.$set(this.lianList[A], 'jieCount', this.jieCount)
          this.zuzhi = 0
          this.jieCount = 0
          this.heightData.OrgName = window.sessionStorage.setItem('OrgName', this.countList.DeployNetCfg.PeerOrgs[0].Name)
        }
      }
    },
    // 增加一行
    removeDomain () {
      this.addForm.DeployCfg.DeployNetCfg.PeerOrgs.push({
        Specs: [{ Hostname: '' }],
        Name: '',
        Domain: '',
        DeployNode: '172-16-254-33',
        Users: {
          Count: 1
        }
      })
    },
    // 删除一行
    deleteItem (index) {
      if (this.addForm.DeployCfg.DeployNetCfg.PeerOrgs.length === 1) return
      this.addForm.DeployCfg.DeployNetCfg.PeerOrgs.splice(index, 1)
    },
    async getClur () {
      const { data: res } = await this.$http.get('clusters')
      console.log('得到集群数据id', res)
      this.cluserData = res.data
      console.log('QQQQQ', this.cluserData)
      for (var i = 0; i < this.cluserData.length; i++) {
        this.clusId = this.cluserData[i].ClusterId
        this.clusId = [...this.cluId, this.cluserData[i].ClusterId]
      }
    },
    // 新建链
    async subCluster () {
      this.$refs.addFormRef.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        // this.addForm.Cert = this.crt // 证书重新赋值
        // this.addForm.Key = this.key // 秘钥重新赋值
        console.log('预校验', this.addForm)
        for (var i = 0; i < this.addForm.DeployCfg.DeployNetCfg.PeerOrgs.length; i++) {
          this.orgDo = this.addForm.DeployCfg.DeployNetCfg.PeerOrgs[i].Name
          console.log(this.orgDo)
          this.addForm.DeployCfg.DeployNetCfg.PeerOrgs[i].Domain = `${this.orgDo}.fabric.baas.xyz`
        }
        // console.log('提交前的数据', this.addForm)
        const { data: res } = await this.$http.post('deploy', this.addForm)
        // console.log('创建区块链数据', res)
        if (res.code !== 0) return this.$message.error('创建区块链失败')
        this.$message.success('创建区块链成功')
        this.paramsChanal.BlockChainId = res.data.BlockChainId
        // this.paramsChanal.BlockChainId = 'DbMh4pVpzoWfCsTyu1dlmF2DyAkGwXGw'
        this.orgChannel.BlockChainId = res.data.BlockChainId
        // this.orgChannel.BlockChainId = 'DbMh4pVpzoWfCsTyu1dlmF2DyAkGwXGw'
        this.paramsChanal.OrgName = this.addForm.DeployCfg.DeployNetCfg.PeerOrgs[0].Name
        const { data: ress } = await this.$http.post('orgcreatechannel', this.paramsChanal)
        if (ress.code !== 0) return this.$message.error('创建通道失败')
        for (var j = 0; j < this.addForm.DeployCfg.DeployNetCfg.PeerOrgs.length; j++) {
          this.org = this.addForm.DeployCfg.DeployNetCfg.PeerOrgs[j].Name
          // console.log(this.org)
          this.orgChannel.OrgName = this.org
          const { data: chunk } = await this.$http.post('orgjoinchannel', this.orgChannel)
          if (chunk.code !== 0) return this.$message.error('加入通道失败')
          this.$message.success('加入通道成功')
        }
        this.addClose = false
        // 隐藏添加用户的对话框
        this.getLianList()
        // 置空数据
        this.$refs.addFormRef.resetFields()
      })
    },
    // 删除链
    async removeLian (name, id) {
      const delForm = {
        BlockChainName: name,
        ClusterId: id
      }
      const isConfird = await this.$confirm(`删除后,『${name}』将停止运行并永久删除?`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      console.log(isConfird)
      if (isConfird !== 'confirm') {
        return this.$message.info('已经取消删除')
      }
      const { data: res } = await this.$http.post('delete', delForm)
      console.log('删除成功', res)
      if (res.code !== 0) return this.$message.error('删除链失败')
      this.$message.success('删除链成功')
      this.getLianList()
    }
  }
}
</script>
<style lang="less" scoped>
.create_contain {
  position: relative;
  .el-dialog {
    .el-dialog__footer {
      text-align: center !important;
    }
  }
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
  .lian {
    width: 741px;
    height: 175px;
    background-color: #f6f8fb;
    border-radius: 12px;
    border: solid 1px #c8c9cc;
    margin-top: 20px;
    .lian_containt {
      padding: 17px;
      .lian_wrap {
        .lian_title {
          display: flex;
          justify-content: space-between;
          color: #717580;
          margin-bottom: 5px;
          .title {
            font-size: 20px;
            color: #1f232d;
          }
          .more {
            font-size: 14px;
            color: #0e85ff;
            cursor: pointer;
          }
        }
        .version {
          color: #717580;
        }
      }
    }
    .count_wrap {
      display: flex;
      margin-top: 32px;
      span {
        color: #717580;
      }
      .zuzhiCount {
        font-size: 30px;
        color: #1f232d;
      }
      div {
        width: 25%;
      }
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
  .el-form-item__content {
    margin-left: 115px;
  }
  .pointAll {
    display: flex;
    .plug {
      margin-bottom: 20px;
      width: 40%;
      margin-right: 10px;
    }
    .plugPoint {
      width: 90%;
    }
  }
  .info {
    margin-top: 10px;
    font-size: 16px;
  }
  .addPlug {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 110px;
    height: 28px;
    background-color: #ebf2ff;
    border-radius: 14px;
    border: solid 1px #2e82ff;
    cursor: pointer;
    margin-left: 106px;
    margin-top: -3px;
    margin-bottom: 10px;
  }
  .anothor {
    margin-top: -22px !important;
  }
}
</style>
