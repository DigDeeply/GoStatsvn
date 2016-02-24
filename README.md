# GoStatsvn
a svn stat tool written by Go.  
用GO写的用来统计每个人的代码提交行数的工具。  
当前是最简版本，日后会持续更新，能更细粒度的统计数据，以及统计每个人每天什么时间点提交代码，每周哪几天提交代码.  

start@2015-03-16 , Go Go Go.


#usage:
用法：
first you need to generate a svn log with xml format, you can also only dump a part of the svn log with -r param.  
首先你需要导出一份xml格式的svn日志,你也可以使用-r参数来限定导出的日志数,避免统计过多.  
then run the GoStatsvn with -f the svn log file, -d the svn working directory and -t the hightcharts template path, it's the gostatsvn.html file  in the project root.  
然后运行编译好的GoStatsvn,使用-f参数指定svn日志文件的位置，-d参数指定svn的开发路径, -t参数指定画图的模版文件路径，模版文件是项目根目录下的gostatsvn.html文件。  
<pre><code>
svn log -v --xml  > svnlog.xml
./GoStatsvn -f svnlog.xml -d workingDirectory -t hightchartsTemplateFilePath
</code></pre>

#demo
![demo](http://p9.qhimg.com/d/inn/6379de03/GoStatsvn-stat-result.png "a GoStatsvn result demo")
#chart
![chart](http://p0.qhimg.com/d/inn/c776308a/qqqq.png "a GoStatsvn result chart demo")
