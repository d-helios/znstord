package znstor

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	API_BASE_PATH = "/api/v1/storage"
	// Domains && Pools && Projects
	DOMAIN_BASE_PATH  = API_BASE_PATH + "/domains/{domain}"
	POOL_BASE_PATH    = DOMAIN_BASE_PATH + "/pools"
	PROJECT_BASE_PATH = POOL_BASE_PATH + "/{pool}/projects"

	// Filesystem
	FILESYSTEM_BASE_PATH          = PROJECT_BASE_PATH + "/{project}/filesystems"
	FILESYSTEM_SNAPSHOT_BASE_PATH = FILESYSTEM_BASE_PATH + "/{filesystem}/snapshots"

	// Volume
	VOLUME_BASE_PATH          = PROJECT_BASE_PATH + "/{project}/volumes"
	VOLUME_SNAPSHOT_BASE_PATH = VOLUME_BASE_PATH + "/{volume}/snapshots"

	// Hosts
	HOST_BASE_PATH = API_BASE_PATH + "/hosts"

	// Targets
	TARGET_BASE_PATH = API_BASE_PATH + "/targets"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(logOutput io.Writer, login, password string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		log.Println("LOAD_ROUTE: ", route)
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Wrapper(handler, route.Name, logOutput, login, password)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{

	/*
		Project Routes
	*/

	Route{
		"ListProjects",
		"GET",
		PROJECT_BASE_PATH,
		HandlerGetProjectList,
	},
	Route{
		"GetProject",
		"GET",
		PROJECT_BASE_PATH + "/{project}",
		HandlerGetProject,
	},
	Route{
		"ProjectExists",
		"GET",
		PROJECT_BASE_PATH + "/{project}/exists",
		HandlerProjectExists,
	},
	Route{
		"CreateProject",
		"POST",
		PROJECT_BASE_PATH + "/{project}",
		HandlerCreateProject,
	},
	Route{
		"ModifyProject",
		"PUT",
		PROJECT_BASE_PATH + "/{project}",
		HandlerModifyProject,
	},
	Route{
		"DestroyProject",
		"DELETE",
		PROJECT_BASE_PATH + "/{project}",
		HandlerDestroyProject,
	},
	Route{
		"ForceDestroyProject",
		"DELETE",
		PROJECT_BASE_PATH + "/{project}/force",
		HandlerForceDestroyProject,
	},

	/*
		Volume Routes
	*/
	Route{
		"ListVolumes",
		"GET",
		VOLUME_BASE_PATH,
		HandlerGetVolumeList,
	},
	Route{
		"GetVolume",
		"GET",
		VOLUME_BASE_PATH + "/{volume}",
		HandlerGetVolume,
	},
	Route{
		"CreateVolume",
		"POST",
		VOLUME_BASE_PATH,
		HandlerCreateVolume,
	},
	Route{
		"DestroyVolume",
		"DELETE",
		VOLUME_BASE_PATH + "/{volume}",
		HandlerDestroyVolume,
	},
	Route{
		"ResizeVolume",
		"PUT",
		VOLUME_BASE_PATH + "/{volume}/resize",
		HandlerResizeVolume,
	},
	Route{
		"SetVolumeCompression",
		"PUT",
		VOLUME_BASE_PATH + "/{volume}/compression/{compression}",
		HandlerSetVolumeCompression,
	},
	Route{
		"CreateVolumeSnapshot",
		"POST",
		VOLUME_SNAPSHOT_BASE_PATH + "/{snapshot}",
		HandlerCreateVolumeSnapshot,
	},
	Route{
		"ListVolumeSnapshots",
		"GET",
		VOLUME_SNAPSHOT_BASE_PATH,
		HandlerGetVolumeSnapshotList,
	},
	Route{
		"GetVolumeSnapshot",
		"GET",
		VOLUME_SNAPSHOT_BASE_PATH + "/{snapshot}",
		HandlerGetVolumeSnapshot,
	},
	Route{
		"RollbackVolumeSnapshot",
		"PUT",
		VOLUME_SNAPSHOT_BASE_PATH + "/{snapshot}/rollback",
		HandlerRollbackVolumeSnapshot,
	},
	Route{
		"DestroyVolumeSnapshot",
		"DELETE",
		VOLUME_SNAPSHOT_BASE_PATH + "/{snapshot}",
		HandlerDestroyVolumeSnapshot,
	},
	Route{
		"CloneFromSnapshot",
		"POST",
		VOLUME_SNAPSHOT_BASE_PATH + "/{snapshot}/clone",
		HandlerCloneVolumeFromSnapshot,
	},
	Route{
		"ExportVolume",
		"PUT",
		VOLUME_BASE_PATH + "/{volume}/export",
		HandlerExportVolume,
	},
	Route{
		"GetExports",
		"GET",
		VOLUME_BASE_PATH + "/{volume}/exports",
		HandlerGetVolumeExports,
	},
	Route{
		"UnExportVolume",
		"PUT",
		VOLUME_BASE_PATH + "/{volume}/unexport",
		HandlerUnExportVolume,
	},
	Route{
		"VolumeJobStatus",
		"GET",
		VOLUME_BASE_PATH + "/job/{uuid}",
		HandlerGetVolumeJobStatus,
	},

	/*
		HostGroup Routes
	*/
	Route{
		"ListHostGroups",
		"GET",
		HOST_BASE_PATH,
		HandlerGetHostGroupList,
	},
	Route{
		"GetHostGroup",
		"GET",
		HOST_BASE_PATH + "/{hostgroup}",
		HandlerGetHostGroup,
	},
	Route{
		"CreateHostGroup",
		"POST",
		HOST_BASE_PATH + "/{hostgroup}",
		HandlerCreateHostGroup,
	},
	Route{
		"AddHostGroupMember",
		"PUT",
		HOST_BASE_PATH + "/{hostgroup}/add/{member}",
		HandlerAddHostGroupMember,
	},
	Route{
		"RemoveHostGroupMember",
		"PUT",
		HOST_BASE_PATH + "/{hostgroup}/remove/{member}",
		HandlerRemoveHostGroupMember,
	},
	Route{
		"DeleteHostGroup",
		"DELETE",
		HOST_BASE_PATH + "/{hostgroup}",
		HandlerDeleteHostGroup,
	},
	Route{
		"AddMultiGroupMember",
		"PUT",
		HOST_BASE_PATH + "/{hostgroup}/add/{member}/force",
		HandlerAddMultiGroupMember,
	},

	/*
		TargetGroup Routes
	*/
	Route{
		"ListTargetGroups",
		"GET",
		TARGET_BASE_PATH + "/tg",
		HandlerGetTargetGroupList,
	},
	Route{
		"GetTargetGroup",
		"GET",
		TARGET_BASE_PATH + "/tg/{targetgroup}",
		HandlerGetTargetGroup,
	},
	Route{
		"CreateTargetGroup",
		"POST",
		TARGET_BASE_PATH + "/tg/{targetgroup}",
		HandlerCreateTargetGroup,
	},
	Route{
		"AddTargetGroupMember",
		"PUT",
		TARGET_BASE_PATH + "/tg/{targetgroup}/add/{member}",
		HandlerAddTargetGroupMember,
	},
	Route{
		"RemoveTargetGroupMember",
		"PUT",
		TARGET_BASE_PATH + "/tg/{targetgroup}/remove/{member}",
		HandlerRemoveTargetGroupMember,
	},
	Route{
		"DeleteTargetGroupMember",
		"DELETE",
		TARGET_BASE_PATH + "/tg/{targetgroup}",
		HandlerDeleteTargetGroup,
	},

	/*
		TargetPortGroup Routes
	*/
	Route{
		"CreateTargetPortGroup",
		"POST",
		TARGET_BASE_PATH + "/tpg/{target_port_group}",
		HandlerCreateTargetPortGroup,
	},
	Route{
		"DeleteTargetPortGroup",
		"DELETE",
		TARGET_BASE_PATH + "/tpg/{target_port_group}",
		HandlerDeleteTargetPortGroup,
	},
	Route{
		"ListTargetPortGroup",
		"GET",
		TARGET_BASE_PATH + "/tpg",
		HandlerGetTargetPortGroupList,
	},
	Route{
		"GetTargetPortGroup",
		"GET",
		TARGET_BASE_PATH + "/tpg/{target_port_group}",
		HandlerGetTargetPortGroup,
	},

	/*
		Target Routes
	*/
	Route{
		"CreateTarget",
		"POST",
		TARGET_BASE_PATH + "/{target}",
		HandlerCreateTarget,
	},
	Route{
		"ListTargets",
		"GET",
		TARGET_BASE_PATH,
		HandlerGetTargetList,
	},
	Route{
		"GetTarget",
		"GET",
		TARGET_BASE_PATH + "/{target}",
		HandlerGetTarget,
	},
	Route{
		"DeleteTarget",
		"DELETE",
		TARGET_BASE_PATH + "/{target}",
		HandlerDeleteTarget,
	},
}
