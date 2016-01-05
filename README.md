# mclink

Rest service written in Go that add/remove minecraft link to/from PC desktop. It also can stop the minecraft application. Tested on a PC running windows 7.

## Info
I'm tired to fight to my son when he has to stop playing Minecraft. I decided to create a service that allow me to control mincraft link on my desktop PC and to stop the minecraft application. It works for me since my kids know to start applications on windows from links on desktop only.

At this moment these are the features I have implemented:

 * add a minecraft link to the desktop, for all users;
 * delete the minectaft link from desktop;
 * stop the minecraft windows process;
 * schedule a stop after some 'minutes' (this will delete the minecraft shortcut on desktop too);
 * warning popup displayed 1 minute before scheduled stop;
  
## Technical details

 * written in Go
 * installed on windows 7 as service using nssm application.
 * application is accessible in home network (I use my mobile phone in order to control it)
 * it uses 'TASKKILL' windows tool to stop the minecraft process
 * it uses 'WMIC PROCESS' in order to determine PID of minecraft process
 * minecraft runs in a java wirtual machine (javaw.exe windows process), therefore I had to determine jvm's PID
 * minecraft has two processes (JVMs): one is the launcher and the other is the main game.
 * minecraft.lnk should exist in minecraft home folder and points to 'minecraft_launcher.exe'
 * this application doesn't create the link but copy it from minecraft home folder
 * *minecraft_home* environment variable should exist and points to the folder minecraft is installed
 * *mclink_log* environment variable is good to be, pointing to a file where this application write logs
 * the appication port is hardcoded to 8123
 * nssm service application should be started with an account with administrativ rights
 * windows firewall should have an inboud rule for 8123 port
  
## Installation steps

Basically, I recommend creating a 'service' folder containing mclink.exe, nssm.exe and mclink.log files. It's good to have all these 3 in the same place.

 * build the application from sources to a file, mclink.exe (*go build*)
 * create a 'service' folder somewhere in the filesystem and move mclink.exe into it
 * create *minecraft_home* environment variable pointing to minecraft home folder
 * create *mclink_log* environment variable pointing to a file in 'service' folder 
 * create the *minecraft.lnk* shortcut for *minecraft_launcher.exe* file in the home folder of minecraft
 * download nssm, copy it into the 'service' folder and use it to install mclink.exe as a service (*nssm install mclink*)
 * go to windows firewall -> inboud rules -> create new rule for 8123 port
 * go to windows services (services.msc as admin) and start mclink service with an admin account (it should successfuly start and you can find a logged line in the mclink.log file in 'service' folder)
 * open a browser locally and test the application accessing http://localhost:8123
 * find the computer IP (*ipconfig*), ex: 192.168.0.10
  
After instalation, this service will start nomatter which user is logged in at that moment on windows.

## Access it from Android phone

 * open a browser and access the application based on IP determined above, ex: http://192.168.0.10:8123 -> you should see the GUI of the application with serveral buttons (Add link, Delete link, Stop, Schedule)  
 * save a bookmark for this page
 * go to phone Android Home and create a shorcut based on previous saved bookmark. (long touch on home screen -> select Widgets -> select a bookmark)

