# 一个简单cron定时任务程序
<pre>
c := cron.New()
c.Add(time.Second * 6, cron.FuncJob(func() {
    fmt.Println(6, time.Now())
}))
c.Start()
</pre>