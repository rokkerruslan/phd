<template>
  <div class="container is-fluid">
    <div class="notification">
      This container is <strong>centered</strong> on desktop.
    </div>

    <div class="column" v-for="event of events">
      <div class="post-module" @mouseover="event.isHover = true" @mouseleave="event.isHover = false" :class="{ hover: event.isHover }">
        <div class="thumbnail">
          <div class="date">
            <div class="day">27</div>
            <div class="month">Мар</div>
          </div>
          <img src="https://s3-us-west-2.amazonaws.com/s.cdpn.io/169963/photo-1429043794791-eb8f26f44081.jpeg" alt=""/>
        </div>
        <div class="post-content">
          <div class="category">Фото</div>
          <h1 class="title" v-text="event.Name"></h1>
          <h2 class="sub_title">The city that never sleeps.</h2>
          <p class="description" v-text="event.Description"></p>
          <div class="post-meta">
            <span class="timestamp"><i class="fa fa-clock-o"></i>6 mins ago</span>
            <span class="comments"><i class="fa fa-comments"></i><a href="#">39 comments</a></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
    const host = "http://localhost:3000/api/v1/events";

    export default {
        name: 'Events',
        data() {
            return {
                hover: false,
                events: []
            }
        },
        async mounted() {
            const data = await fetch(host);
            if (data.ok) {
                const events = await data.json();
                events.forEach(e => e.isHover = false);
                this.events = events
            } else {
                console.log(data.error())
            }
        }
    }
</script>

<style scoped>
  .container {
    font-size: 14px;
  }

  .post-module {
    position: relative;
    z-index: 1;
    display: block;
    background: #FFFFFF;
    min-width: 270px;
    height: 470px;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.15);
    transition: all 0.3s linear 0s;
  }

  .hover {
    box-shadow: 0 1px 35px 0 rgba(0, 0, 0, 0.3);
  }

  .post-module:hover .thumbnail img, .hover .thumbnail img {
    transform: scale(1.1);
    opacity: 0.6;
  }

  .post-module .thumbnail {
    background: #000000;
    height: 400px;
    overflow: hidden;
  }

  .post-module .thumbnail .date {
    position: absolute;
    top: 20px;
    right: 20px;
    z-index: 1;
    background: #db4200;
    width: 55px;
    height: 55px;
    padding: 12.5px 0;
    border-radius: 100%;
    color: #FFFFFF;
    font-weight: 700;
    text-align: center;
    box-sizing: border-box;
    border: 2px solid white;
  }

  .post-module .thumbnail .date .day {
    font-size: 18px;
    top: -9px;
    position: relative;
  }

  .post-module .thumbnail .date .month {
    font-size: 12px;
    text-transform: uppercase;
    top: -13px;
    position: relative;
  }

  .post-module .thumbnail img {
    display: block;
    width: 120%;
    transition: all 0.3s linear 0s;
  }

  .post-module .post-content {
    position: absolute;
    bottom: 0;
    background: #FFFFFF;
    width: 100%;
    padding: 30px;
    box-sizing: border-box;
    transition: all 0.3s cubic-bezier(0.37, 0.75, 0.61, 1.05) 0s;
  }

  .post-module .post-content .category {
    position: absolute;
    top: -34px;
    left: 0;
    background: #e74c3c;
    padding: 10px 15px;
    color: #FFFFFF;
    font-size: 14px;
    font-weight: 600;
    text-transform: uppercase;
  }

  .post-module .post-content .title {
    margin: 0;
    padding: 0 0 10px;
    color: #333333;
    font-size: 26px;
    font-weight: 700;
  }

  .post-module .post-content .sub_title {
    margin: 0;
    padding: 0 0 20px;
    color: #e74c3c;
    font-size: 20px;
    font-weight: 400;
  }

  .post-module .post-content .description {
    display: none;
    color: #666666;
    font-size: 14px;
    line-height: 1.8em;
  }

  .post-module .post-content .post-meta {
    margin: 30px 0 0;
    color: #999999;
  }

  .post-module .post-content .post-meta .timestamp {
    margin: 0 16px 0 0;
  }

  .post-module .post-content .post-meta a {
    color: #999999;
    text-decoration: none;
  }

  .hover .post-content .description {
    display: block !important;
    height: auto !important;
    opacity: 1 !important;
  }

  .container .column {
    max-width: 400px;
    min-width: 320px;
    width: 50%;
    padding: 0 25px;
    box-sizing: border-box;
    float: left;
  }

</style>