<template>
  <div class="login_container">
    <div class="login_left">
      <div class="login_log">
        <img src="../assets/image/logo1.png" />
      </div>
      <div class="avatar_title">
        <div>你好！</div>
        <div>欢迎登录链飞控制台</div>
      </div>
      <div class="login_box">
        <!-- 表单额数据绑定:v-model  表单校验 :rules: 怎么才能让表单的校验生效 prop=""
        表单的引用ref=""：可以直接调用表单的方法-->
        <el-form
          label-width="70px"
          class="formD"
          :model="LoginData"
          :rules="LoginRules"
          ref="resetForm"
        >
          <!-- 用户名 -->
          <el-form-item label="邮箱">
            <el-input prop="Mail" placeholder="请输入邮箱" v-model="LoginData.Mail">
              <!-- <template slot="prepend">邮箱</template> -->
            </el-input>
          </el-form-item>
          <!-- 密码 -->
          <el-form-item label="密码">
            <el-input prop="Password" placeholder="请输入" v-model="LoginData.Password"></el-input>
          </el-form-item>
          <!-- <div class="forgetPsd">忘记密码？</div> -->
          <!-- 登录重置 -->
          <!-- <el-form-item class="btns"> -->
          <div class="contain_register">
            <el-button type="primary" class="sub-btn" @click="login()">登录</el-button>
            <div class="register">
              还没有账号？
              <span class="rightRegist" @click="goRegist()">立即注册</span>
              <!-- <a class="rightRegist" href="https://wingchain.cn/consulting_business?consultType=1" target="_blank">申请试用</a> -->
            </div>
          </div>
          <!-- </el-form-item> -->
        </el-form>
      </div>
    </div>
    <div class="login_right" @click="go">
      <img src="../assets/image/bg.png" alt />
    </div>
  </div>
</template>
<script>
export default {
  data () {
    // 校验邮箱
    var checkEmail = (rule, value, callback) => {
      const mailReg = /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+/
      if (!value) {
        return callback(new Error('邮箱不能为空'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('请输入正确的邮箱格式'))
        }
      }, 100)
    }
    return {
      LoginData: {
        // Mail: '',
        // Password: ''
        Mail: 'a@qq.com',
        Password: 'aaaaaaa'
      },
      hre: 'https://wingchain.cn/consulting_business?consultType=1',
      LoginRules: {
        Mail: [{ required: true, validator: checkEmail, trigger: 'blur' }],
        Password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
          // { min: 3, max: 9, message: '长度在 3 到 5 个字符', trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    goLian (e) {
      console.log(e)
      window.localtion.href = 'https://www.baidu.com'
      // const { data: res } = await this.$http.post('https://wingchain.cn/consulting_business?consultType=1')
      // console.log('跳转链接', res)
    },
    // 跳转至注册页面
    goRegist () {
      console.log('aaa')
      this.$router.push('/registered')
    },
    //   重置数据
    resetForm () {
      // 调用的是方法
      this.$refs.resetForm.resetFields()
    },
    async login () {
      // 登录验证
      this.$refs.resetForm.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return // false不合法，直接返回
        const { data: res } = await this.$http.post('login', this.LoginData)
        console.log(res)
        console.log(this.LoginData.Mail)
        window.sessionStorage.setItem('Mail', this.LoginData.Mail)
        if (res.code !== 0) return this.$message.success('登录失败,用户名或密码错误')
        this.$message.success('登录成功')
        /*
        1: 将登录成果的token，保存在客户端的sessionStroge中
           1.1项目中除了登录之外的api,必须在登录之后才能访问
           1.2 token只应在当前打开的网站生效，所以将token保存在sessionStroge中
        2：通过编程式导航到后台主页，/home
        */
        // window.sessionStorage.setItem('token', res.data.token)
        this.$router.push('/home')
      })
    },
    go () {
      this.$router.push('/home')
    }
  }
}
</script>
<style lang="less" scoped>
.login_container {
  background: #fff;
  height: 100%;
  position: relative;
  display: flex;
  justify-content: space-between;
  .login_left {
    .login_log {
      position: absolute;
      top: 24px;
      left: 62px;
      width: 241px;
      height: 62px;
      img {
        width: 100%;
        height: 100%;
      }
    }
    .avatar_title {
      position: absolute;
      top: 157px;
      left: 71px;
      div {
        font-family: PingFangSC-Semibold;
        font-size: 32px;
        font-weight: normal;
        font-stretch: normal;
        color: #1f232d;
      }
    }
    .login_box {
      width: 325px;
      height: 248px;
      background: #fff;
      border-radius: 3px;
      position: absolute;
      top: 265px;
      left: 31px;
      .avatar_box {
        width: 130px;
        height: 130px;
        border: 1px solid #eee;
        border-radius: 50%;
        padding: 10px;
        box-shadow: 0 0 10px #eee;
        position: absolute;
        left: 50%;
        transform: translate(-50%, -50%);
        img {
          width: 100%;
          height: 100%;
          border-radius: 50%;
        }
      }
      .formD {
        position: absolute;
        bottom: 0;
        width: 100%;
        box-sizing: border-box;
        .el-input-group__prepend {
          background: #fff !important;
          border: none !important;
        }
        .forgetPsd {
          width: 70px;
          height: 12px;
          float: right;
          font-family: PingFangSC-Medium;
          font-size: 12px;
          font-weight: normal;
          font-stretch: normal;
          line-height: 10px;
          letter-spacing: 0px;
          color: #2e82ff;
          margin-top: -13px;
          cursor: pointer;
        }
        .contain_register {
          .el-button--primary.sub-btn {
            width: 77%;
            // margin-top: 13px;
            margin-top: 6px;
            margin-left: 73px;
          }
          .register {
            width: 176px;
            height: 13px;
            margin: auto;
            font-family: PingFangSC-Medium;
            font-size: 14px;
            font-weight: normal;
            font-stretch: normal;
            line-height: 10px;
            letter-spacing: 0px;
            color: #1f232d;
            margin-top: 20px;
            margin-left: 107px;
            .rightRegist {
              font-family: PingFangSC-Medium;
              font-size: 14px;
              font-weight: normal;
              letter-spacing: 0px;
              color: #2e82ff;
              cursor: pointer;
            }
            .rightRegist {
              font-family: PingFangSC-Medium;
              font-size: 14px;
              font-weight: normal;
              letter-spacing: 0px;
              color: #2e82ff;
              cursor: pointer;
              text-decoration-line: none;
            }
          }
        }
      }
    }
  }
  .login_right {
    width: 1126px;
    height: 650px;
    margin-top: 63px;
    img {
      width: 100%;
      height: 100%;
    }
  }
  // /deep/ .el-input {
  // .el-input {
  //   width:60%
  // }
}
</style>
