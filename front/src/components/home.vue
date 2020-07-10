<template>
  <!-- home页的主体布局 -->
  <el-container class="home-container">
    <!-- 头部区域 -->
    <el-container>
      <!-- 侧边栏 -->
      <el-aside :width="isCollapse ? '64px' : '200px'">
        <div class="logo">
          <img src="../assets/image/logo2.png" alt />
        </div>
        <div>
          <!-- <div class="logo_title" v-if="$store.state.isShow" @click="goMy()">{{$store.state.title}} -->
          <!-- <div class="logo_title" v-if="this.isShow" @click="goMy()"> -->
          <div class="logo_title" v-if="$store.state.isShow" @click="goMy()">
            {{$store.state.title}}
            <i class="el-icon-arrow-right"></i>
          </div>
          <div class="logo_title" v-else>
            我的联盟
            <i class="el-icon-arrow-right"></i>
          </div>
        </div>
        <!-- background-color="#1f232d" -->
        <!-- 菜单栏分为两级，并且可以折叠 -->
        <div v-if="$store.state.isShow">
        <!-- <div v-if="this.isShow"> -->
          <el-menu
            background-color="#1f232d"
            active-text-color="#fff"
            text-color="#fff"
            unique-opened
            :collapse="isCollapse"
            :collapse-transition="false"
            router
            :default-active="activePath"
          >
            <!-- submun :index="" string  + ‘’ 转化为字符串 -->
            <el-menu-item index="deploy">
              <!-- <i class="el-icon-menu"></i> -->
              <i class="iconfont" style="font-size:13px">&#xe689;</i>
              <span slot="title">区块链管理</span>
            </el-menu-item>
            <el-menu-item index="treaty">
              <i class="iconfont">&#xe687;</i>
              <span slot="title">合约管理</span>
            </el-menu-item>
            <el-menu-item index="user">
              <i class="iconfont">&#xe686;</i>
              <span slot="title">联盟管理</span>
            </el-menu-item>
            <el-menu-item index="cluster">
              <i class="iconfont">&#xe688;</i>
              <span slot="title">集群管理</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-aside>
      <!-- 主体区 -->
      <el-main>
        <!-- 路由展位符: 在home页面的路由展位，指派children -->
        <router-view></router-view>
      </el-main>
      <div class="emai">
        <el-dropdown class="user-name" split-button trigger="click" @command="handleCommand">
          <span class="el-dropdown-link">
            {{this.Mail}}
          </span>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="loginout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>

      <!-- <div class="emai">{{this.Mail}}
      <i class="el-icon-caret-bottom" type="primary" @click="loginOut"></i>
      </div>-->
    </el-container>
  </el-container>
</template>
<script>
export default {
  data () {
    return {
      title: this.$title,
      menuList: [],
      Mail: '', //   取来的邮箱
      iconObj: {
        125: 'iconfont icon-users',
        103: 'iconfont icon-tijikongjian',
        101: 'iconfont icon-3702mima',
        102: 'iconfont icon-danju',
        145: 'iconfont icon-baobiao'
      },
      isCollapse: false,
      activePath: '',
      AllianceId: '',
      AllianceName: '',
      // isShow: '',
      Title: ''
    }
  },
  created () {
    this.Mail = window.sessionStorage.getItem('Mail')
    this.AllianceId = window.localStorage.getItem('AllianceId')
    this.AllianceName = window.localStorage.getItem('AllianceName')
    // this.Title = window.localStorage.getItem('title')
    // this.isShow = window.localStorage.getItem('isShow')
    console.log('看看一进来是什么', this.isShow)
  },
  mounted () {
  },
  methods: {
    goMy () {
      this.$router.push('/create')
    },
    roggleButton () {
      console.log('123')
      this.isCollapse = !this.isCollapse
    },
    /*
        退出功能
        只需要清除token，跳转到登录页。这样后续的请求就不会携带token，必须生成一个新的token才能访问页面
*/
    backHome () {
      window.sessionStorage.clear()
      this.$router.push('/login')
    },
    handleCommand (command) {
      if (command === 'loginout') {
        // 清除用户信息
        this.$store.commit('changeshow')
        this.$router.push('/login')
        // window.localStorage.setItem('isShow', false)
      }
    }

  }
}
</script>
<style lang="less" scoped>
.home-container {
  height: 100%;
  .emai {
    cursor: pointer;
    width: 160px;
    position: absolute;
    top: 13px;
    right: 23px;
  }
  /deep/ .el-menu-item:hover {
    background: #2f333e !important;
  }
}
.el-header {
  background: #373d41;
  display: flex;
  justify-content: space-between;
  color: #fff;
  padding-left: 0;
  align-items: center;
  > div {
    display: flex;
    align-items: center;
    > span {
      margin-left: 30px;
    }
  }
}
.el-aside {
  background: #1f232d;
  .logo {
    width: 175px;
    height: 34px;
    padding-top: 9px;
    padding-left: 3px;
    img {
      width: 100%;
    }
  }
  .logo_title {
    font-size: 18px;
    color: #ffffff;
    // padding-left: 19px;
    padding-left: 34px;
    padding-bottom: 10px;
    margin-top: 30px;
    cursor: pointer;
  }
  .el-menu {
    border-right: none;
    font-size: 16px;
    color: #ffffff;
    background: #2f333e;
    .el-menu-item {
      background: #2f333e;
    }
  }
}
.el-main {
  background: #fff;
  display: flex;
  // justify-content: space-between;
  margin-bottom: 15px;
}
.iconfont {
  margin-right: 10px;
}
.roggleButton {
  color: #fff;
  text-align: center;
  background: #4a5064;
  font-size: 10px;
  line-height: 24px;
  letter-spacing: 0.2em;
  cursor: pointer;
}
</style>
