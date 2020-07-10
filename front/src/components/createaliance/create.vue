<template>
  <div class="create_contain">
    <div class="lince_title">联盟管理</div>
    <div class="my_create">
      <div class="line"></div>我创建的联盟
    </div>
    <div class="add_all">
      <div class="add_liance">
        <div class="add_box" @click="addClose = true">
          <div class="add_line">+</div>
          <div class="add_title">新建联盟</div>
        </div>
      </div>
      <div class="linceList" v-for="(item, index) in allinceList" :key="index">
        <div class="del_llince" @click="del_allince(item.Id, item.Name)">...</div>
        <div class="add_box" @click="lian_Main(item.Id, item.Name)">
          <div class="add_line">{{item.Name}}</div>
          <div class="add_title">
            区块链数量：
            <!-- <span class="add_desc" >{{item.len}}</span> -->
            <span class="add_desc">{{item.len}}</span>
          </div>
          <div class="add_title" style="margin-top:8px">
            描述：
            <span class="add_desc">{{item.Describe}}</span>
          </div>
        </div>
      </div>
    </div>
    <div class="my_create">
      <div class="line"></div>我加入的联盟
    </div>
    <div class="add_all">
      <div class="add_liance">
        <div class="add_box" @click="joinClose = true">
          <div class="add_line">+</div>
          <div class="add_title">加入联盟</div>
        </div>
      </div>
      <div class="linceList" v-for="(item, index) in joinallinceList" :key="index">
        <div class="del_llince"></div>
        <div class="add_box" @click="lian_Main(item.Id, item.Name)">
          <div class="add_line">{{item.Name}}</div>
          <div class="add_title">
            区块链数量：
            <span class="add_desc">{{item.len}}</span>
          </div>
          <div class="add_title" style="margin-top:8px">
            描述：
            <span class="add_desc">{{item.Describe}}</span>
          </div>
        </div>
      </div>
      <!-- 添加用户对话框  :visible.sync 控制显示隐藏-->
      <div>
        <el-dialog title="新建联盟" :visible.sync="addClose" width="50%" class="add_lince">
          <!-- 内容主体区 -->
          <!-- model数据绑定，验证规则，引用对象 -->
          <el-form :model="addForm" :rules="addFormRules" ref="addFormRef" label-width="100px">
            <el-form-item prop="Name" label="联盟名称">
              <el-input v-model="addForm.Name" maxlength="20" show-word-limit></el-input>
            </el-form-item>
            <el-form-item prop="Describe" label="联盟描述">
              <el-input v-model="addForm.Describe" maxlength="200" show-word-limit></el-input>
            </el-form-item>
          </el-form>
          <!-- 底部 -->
          <span slot="footer" class="dialog-footer">
            <el-button @click="addLince">创建</el-button>
            <el-button type="primary" @click="addClose = false">取消</el-button>
          </span>
        </el-dialog>
      </div>
      <div>
        <el-dialog title="邀请机构" :visible.sync="joinClose" width="50%" class="add_linceVane">
          <!-- 内容主体区 -->
          <!-- model数据绑定，验证规则，引用对象 -->
          <el-form
            :model="joinForm"
            :rules="addVanceFormRules"
            ref="addVanceFormRef"
            label-width="162px"
          >
            <el-form-item prop="Alliance.Id" label="受邀机构电子邮箱">
              <el-input
                v-model="joinForm.Alliance.Id"
                style="margin-left:4px"
                placeholder="请输入你收到的联盟邀请码"
              ></el-input>
            </el-form-item>
          </el-form>
          <!-- 底部 -->
          <span slot="footer" class="dialog-footer">
            <el-button @click="joinLince(joinForm.Alliance.Id)" type="primary">确定</el-button>
            <!-- <el-button type="primary" @click="joinClose = false">取消</el-button> -->
          </span>
        </el-dialog>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  data () {
    return {
      addClose: false, // 新建联盟对话框状态
      joinClose: false, // 加入联盟对话框状态
      allinceList: [], // 联盟列表数据
      joinallinceList: [], // 加入的联盟
      Mail: '', //   取来的邮箱
      Alliance: '',
      AllianceJoin: '',
      surpList: [],
      lianListLen: 0,
      // 添加用户的表单数据
      addForm: {
        Name: '',
        Describe: '',
        // Creator: this.Mail
        Creator: ''
      },
      // 加入联盟的表单数据
      joinForm: {
        Mail: '',
        Alliance: {
          Name: '',
          Describe: '',
          Creator: '',
          Id: ''
        }
      },
      addVanceFormRules: {
        Alliance: {
          Id: [
            {
              required: true,
              message: '请输入受邀机构邮箱',
              trigger: 'blur'
            }
          ]
        }
      },
      // 表单的验证规则
      addFormRules: {
        Name: [
          {
            required: true,
            message: '请输入联盟名称',
            trigger: 'blur'
          },
          {
            max: 20,
            message: '联盟名称的长度在20个字符之间',
            trigger: 'blur'
          }
        ],
        Describe: [
          {
            required: false,
            message: '请输入联盟描述',
            trigger: 'blur'
          },
          {
            max: 200,
            message: '联盟描述的长度在200个字符之间',
            trigger: 'blur'
          }
        ]
      }
    }
  },
  created () {
    this.Mail = window.sessionStorage.getItem('Mail')
    this.getlinceList()
    this.getjoinList()
  },
  methods: {
    // 为了计算链数量
    async getLianList (P) {
      const { data: res } = await this.$http.get(`${this.Alliance}/alliancechains`)
      console.log('联盟之下拉取链列表', res)
      this.lianList = res.data
      console.log('每次的数据', this.lianList)
      // 计算加入的联盟和创建的联盟，里面有多少条“区块链数量”
      if (this.lianList) {
        this.$set(this.allinceList[P], 'len', this.lianList.length)
        this.$set(this.joinallinceList[P], 'len', this.lianList.length)
      }
      console.log('此时创建的联盟列表', this.allinceList)
      console.log('此时加入的联盟列表', this.joinallinceList)
    },
    // 保存联盟id
    lian_Main (id, Name) {
      window.localStorage.setItem('AllianceId', id)
      // this.$title = window.localStorage.setItem('AllianceName', Name)
      window.localStorage.setItem('AllianceName', Name)
      this.$store.commit('changeTitle', Name)
    },
    // 我加入的联盟列表
    async getjoinList () {
      const { data: res } = await this.$http.get(`${this.Mail}/joinedalliances`)
      console.log('加入的联盟列表', res)
      this.joinallinceList = res.data
      if (this.joinallinceList) {
        for (var K = 0; K < this.joinallinceList.length; K++) {
          this.AllianceJoin = this.joinallinceList[K].Id
          this.getLianList(K)
        }
      }
    },
    // 删除联盟
    async del_allince (id, name) {
      const delForm = {
        Mail: this.Mail,
        AllianceId: id
      }
      const isConfird = await this.$confirm(`此操作将永久删除『${name}』, 是否继续`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      console.log(isConfird)
      if (isConfird !== 'confirm') {
        return this.$message.info('已经取消删除')
      }
      const { data: res } = await this.$http.post('deletealliance', delForm)
      console.log('删除成功', res)
      if (res.code !== 0) return this.$message.error('删除联盟失败')
      this.$message.success('删除联盟成功')
      this.getlinceList() // 创建的联盟列表
      this.getjoinList() // 加入的联盟列表
    },
    // 我创建的联盟列表数据
    async getlinceList () {
      const { data: res } = await this.$http.get(`${this.Mail}/createdalliances`)
      console.log('创建的联盟列表数据', res)
      this.allinceList = res.data
      if (this.allinceList) {
        for (var i = 0; i < this.allinceList.length; i++) {
          this.Alliance = this.allinceList[i].Id
          this.getLianList(i)
        }
      }
    },
    // 创建联盟
    async addLince () {
      //   const { data: res } = await this.$http.post('createaliance', this.addForm)
      this.$refs.addFormRef.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        this.addForm.Creator = this.Mail
        console.log('预校验', this.addForm)
        const { data: res } = await this.$http.post('createalliance', this.addForm)
        console.log('新增联盟数据', res)
        if (res.code !== 0) return this.$message.error('创建联盟失败,已存在此联盟')
        this.$message.success('创建联盟成功')
        // 隐藏添加用户的对话框
        this.addClose = false
        // 重新刷新列表获取数据
        this.getlinceList()
        this.getjoinList()
        this.addForm.Name = ''
        this.addForm.Describe = ''
      })
    },
    // 加入联盟
    async joinLince (id) {
      //   const { data: res } = await this.$http.post('createaliance', this.addForm)
      this.$refs.addVanceFormRef.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        this.addForm.Creator = this.Mail
        console.log('预校验', this.addForm)
        // this.joinForm.Alliance.Name =
        const { data: res } = await this.$http.get(`${id}/alliance`)
        console.log('拉取联盟数据', res.data)
        if (res.code !== 0) return this.$message.error('加入联盟失败')
        this.joinForm.Mail = this.Mail
        this.joinForm.Alliance = res.data
        console.log(this.joinForm)
        const { data: resAdd } = await this.$http.post('useraddalliance', this.joinForm)
        console.log('加入联盟数据', resAdd)
        if (resAdd.code !== 0) return this.$message.error('加入联盟失败')
        // 隐藏添加用户的对话框
        this.joinClose = false
        // 重新刷新列表获取数据
        this.getjoinList()
        this.joinForm.Alliance.Id = ''
      })
    }
  }
}
</script>
<style lang="less" scoped>
.create_contain {
  position: relative;
  .add_linceVane {
    /deep/ .el-dialog__body {
      padding-bottom: 0;
    }
  }
  .emai {
    position: absolute;
    top: 0px;
    right: 10px;
  }
  .lince_title {
    font-size: 24px;
    color: #1f232d;
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
  .add_all {
    display: flex;
    .add_liance {
      margin-top: 21px;
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
      margin-left: 10px;
      margin-top: 21px;
      width: 300px;
      height: 133px;
      background-color: #f6f8fb;
      border-radius: 12px;
      border: solid 1px #c8c9cc;
      position: relative;
      .del_llince {
        color: #c8c9cc;
        font-size: 30px;
        margin-left: 115px;
        position: absolute;
        right: 15px;
        top: -16px;
        cursor: pointer;
      }
      .add_box {
        padding-top: 17px;
        padding-left: 15px;
        padding-right: 15px;
        height: 100%;
        font-size: 18px;
        color: #1f232d;
        cursor: pointer;
        .add_line {
          margin-bottom: 10px;
        }
        .add_title {
          font-size: 14px;
          color: #717580;
          margin-top: 23px;
          .add_desc {
            color: #1f232d;
          }
        }
      }
    }
  }
  .el-dialog__header {
    background: red;
  }
}
</style>
