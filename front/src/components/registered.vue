<template>
  <div class="login_container">
    <div class="login_left">
      <div class="login_log">
        <img src="../assets/image/logo1.png" />
      </div>
      <div class="avatar_title">
        <div>你好！</div>
        <div>欢迎注册</div>
      </div>
      <div class="login_box">
        <!-- 表单额数据绑定:v-model  表单校验 :rules: 怎么才能让表单的校验生效 prop=""
        表单的引用ref=""：可以直接调用表单的方法-->
        <el-form
          label-width="0px"
          class="formD"
          :model="LoginData"
          :rules="LoginRules"
          ref="resetForm"
        >
          <!-- 用户名 -->
          <el-form-item prop="Mail">
            <el-input placeholder="请输入邮箱" v-model="LoginData.Mail"></el-input>
          </el-form-item>
          <!-- 密码 -->
          <el-form-item prop="Password">
            <el-input placeholder="设置密码" v-model="LoginData.Password"></el-input>
          </el-form-item>
          <el-form-item prop="repassword">
            <el-input placeholder="确认密码" v-model="LoginData.repassword"></el-input>
          </el-form-item>
          <!-- 登录重置 -->
          <!-- <el-form-item class="btns"> -->
          <el-button type="primary" class="sub-btn" @click="rightRegister()">立即注册</el-button>
          <div class="register">
            已有账号？
            <span class="rightRegist" @click="gologin">立即登录</span>
          </div>
          <!-- </el-form-item> -->
        </el-form>
      </div>
    </div>
    <div class="login_right">
      <img src="../assets/image/bg.png" alt />
    </div>
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
    // 校验密码
    var checkPwd = (rule, value, callback) => {
      // const mailReg = /^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,18}$/ 只能是数字和密码的组合(密码强度)
      const mailReg = /[a-zA-Z0-9-]{6,18}$/ // 可以是数字，字母，数字字母的组合外加-
      if (!value) {
        return callback(new Error('请设置密码'))
      }
      setTimeout(() => {
        if (mailReg.test(value)) {
          callback()
        } else {
          callback(new Error('密码格式不正确，请输入6~18位英文或数字组合'))
        }
      }, 100)
    }
    // 二次校验
    var validatePass = (rule, value, callback) => {
      console.log('密码', this.LoginData.Password)
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.LoginData.Password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    return {
      LoginData: {
        Mail: '',
        Password: '',
        VerifyCode: ''
      },
      LoginRules: {
        Mail: [{ required: true, validator: checkEmail, trigger: 'blur' }],
        Password: [
          { required: true, validator: checkPwd, trigger: 'blur' }
          // { min: 3, max: 9, message: '长度在 3 到 5 个字符', trigger: 'blur' }
        ],
        repassword: [{ required: true, validator: validatePass, trigger: 'blur' }]
      }
    }
  },
  methods: {
    gologin () {
      this.$router.push('/login')
    },
    // async goLian () {
    //   const { data: res } = await this.$http.post('https://wingchain.cn/consulting_business?consultType=1')
    //  console.log('跳转链接', res)
    // },
    async rightRegister () {
      // 登录验证
      this.$refs.resetForm.validate(async valid => {
        //  validate 接收一个参数作为回调函数，成功返回true，返回false
        console.log(valid)
        if (!valid) return
        const { data: res } = await this.$http.post('register', this.LoginData)
        console.log(res)
        if (res.code !== 0) return this.$message.success('注册失败,用户邮箱已存在')
        this.$message.success('注册成功')
        this.$router.push('/login')
      })
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
      top: 163px;
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
      width: 374px;
      height: 248px;
      background: #fff;
      border-radius: 3px;
      position: absolute;
      top: 344px;
      left: 55px;
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
        padding: 0 20px;
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
        .el-button--primary.sub-btn {
          width: 100%;
          margin-top: 15px;
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
          margin-top: 25px;
          .rightRegist {
            font-family: PingFangSC-Medium;
            font-size: 14px;
            font-weight: normal;
            letter-spacing: 0px;
            color: #2e82ff;
            cursor: pointer;
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
}
</style>
