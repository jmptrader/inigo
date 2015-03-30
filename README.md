# Inigo
A lightweight service management tool written in Go.

# Features
* Works independently in user home directory.
* Save current service list to file for easy loading/unloading different sets of services.
* Reboot all services with a single command.

#### Notes:
- All working files reside in the HOME/.inigo directory unless otherwise specified (save/load).
- Current services are automatically saved to default services file (~/.inigo/services) on every add, remove, enable, and disable.

# Examples

#### Start the server
Start up the server to begin managing services. Ideally you should run this in the background.

```
inigo -server
```

#### Shutdown the server
Shutdown will interrupt all services and close the server program.

```
inigo -shutdown
```

#### Add a service
When adding a service, the program name is used as the service name by default. But you can edit the services file to change it.

```
inigo -add ping www.google.com
```

#### Start a service
Services are enabled by default when added. But they are not started automatically until next 'boot'. To manually start a service do:

```
inigo -start ping
```

#### Stop a service
Stopping a service sends an interrupt signal to the process then removes it from the process map.

```
inigo -stop ping
```

#### Remove a service
Removing a service removes it from the current services list (as long as its not running).

```
inigo -remove ping
```

#### Enable service
Enabling a service sets it to run automatically on 'boot'. Services are enabled by default when added.

```
inigo -enable ping
```

#### Disable service
Disabling a service will stop it from running automatically at 'boot'.

```
inigo -disable ping
```

#### Save services to file
Services can be saved to a separate file to create different sets of services.

```
inigo -save ~/path/file
```

#### Unload all services
Unload will interrupt and unload all current processes & services.

```
inigo -unload
```

#### Load services from file
Load the set of services from a specific file.

```
inigo -load ~/path/file
```

#### Reload current services
Reload all current services from default services file.

```
inigo -reload
```

#### Reboot all services
Interrupt all services then boot all currently enabled services.

```
inigo -reboot
```