<template>
  <div id="app">
    <div id="nav">
      <router-link to="/">События</router-link> |
      <router-link to="/about">Как это работает?</router-link> |
      <router-link to="/blog">Блог</router-link> |
      <router-link to="/profile">Аккаунт</router-link>
    </div>
    <router-view/>
  </div>
</template>

<style>
@import "styles/main.css";

#app {
  font-family: 'Oswald', sans-serif;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #db4200;
}
</style>
