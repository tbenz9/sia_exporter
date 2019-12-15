package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/NebulousLabs/Sia/modules"
	"gitlab.com/NebulousLabs/Sia/node/api"
	sia "gitlab.com/NebulousLabs/Sia/node/api/client"
	"gitlab.com/NebulousLabs/errors"
)

var (
	// ErrAPICallNotRecognized is returned by API client calls made to modules that
	// are not yet loaded.
	ErrAPICallNotRecognized = errors.New("API call not recognized")

	// Define the metrics we wish to expose
	// Renter Metrics
	renterModuleLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_module_loaded", Help: "Is the renter module loaded. 0=not loaded.  1=loaded"})
	renterAggregateNumFiles = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_aggregate_num_files", Help: "Shows the number of files uploaded to Sia by the renter"})
	renterAggregateNumStuckChunks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_aggregate_num_stuck_chunks", Help: "The aggregate number of stuck chunks"})
	renterAggregateSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_aggregate_size", Help: "The aggregate size of data stored on Sia"})
	renterMaxHealth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_max_health", Help: "The max health"})
	renterMinRedundancy = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_min_redundancy", Help: "The min redundancy"})
	renterRateLimitDownload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_rate_limit_download", Help: "renter download ratelimit (bytes-per-second)"})
	renterRateLimitUpload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_rate_limit_upload", Help: "renter upload ratelimit (bytes-per-second)"})
	// Contracts
	renterNumActiveContracts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_num_active_contracts", Help: "Number of active contracts"})
	renterNumDisabledContracts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_num_disabled_contracts", Help: "Number of disabled contracts"})
	renterNumRefreshedContracts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_num_refreshed_contracts", Help: "Number of refreshed contracts"})
	renterNumPassiveContracts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_num_passive_contracts", Help: "Number of passive contracts"})
	renterNumExpiredContracts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_num_expired_contracts", Help: "Number of expired contracts"})
	renterNumExpiredRefreshedContracts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_num_expired_refreshed_contracts", Help: "Number of expired refreshed contracts"})
	// Allowance
	renterAllowanceAmount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_amount", Help: "Renter allowance Amount (siacoins)"})
	renterAllowancePeriod = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_period", Help: "Renter allowance period length (blocks)"})
	renterAllowanceRenewWindow = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_renew_window", Help: "Renter allowance renew window (blocks)"})
	renterAllowanceHosts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_hosts", Help: "Renter allowance hosts"})
	renterAllowanceCurrentSpent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_spent", Help: "Amount of allowance in Siacoins spent in the current period"})
	renterAllowanceCurrentUnspent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_unspent", Help: "Unspent amount of allowance in Siacoins in the current period"})
	renterAllowanceCurrentStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_storage", Help: "Amount of allowance in Siacoins spent in the current period on storage"})
	renterAllowanceCurrentUpload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_upload", Help: "Amount of allowance in Siacoins spent in the current period on upload bandwidth"})
	renterAllowanceCurrentDownload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_download", Help: "Amount of allowance in Siacoins spent in the current period on download bandwidth"})
	renterAllowanceCurrentFees = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_fees", Help: "Amount of allowance in Siacoins spent in the current period on fees"})
	renterAllowanceCurrentUnspentAllocated = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_unspent_allocated", Help: "Amount of allocated unspent allowance in Siacoins"})
	renterAllowanceCurrentUnspentUnallocated = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "renter_allowance_current_unspent_unallocated", Help: "Amount of unallocated unspent allowance in Siacoins"})

	// Consensus Metrics
	consensusModuleLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "consensus_module_loaded", Help: "Is the consensus module loaded. 0=not loaded.  1=loaded"})
	consensusSynced = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "consensus_synced", Help: "Consensus sync status, 0=not synced.  1=synced"})
	consensusHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "consensus_height", Help: "Consensus block height"})
	consensusDifficulty = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "consensus_difficulty", Help: "Consensus difficulty"})
	consensusHshrt = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "consensus_hashrate", Help: "Consensus hashrate"})

	// Daemon Metrics
	//	daemonAggregateNumAlerts = promauto.NewGauge(prometheus.GaugeOpts{
	//		Name: "daemon_aggregate_num_alerts", Help: "Total number of daemon Alerts"})
	daemonRateLimitDownload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "global_rate_limit_download", Help: "global download ratelimit (bytes-per-second)"})
	daemonRateLimitUpload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "global_rate_limit_upload", Help: "global upload ratelimit (bytes-per-second)"})

	// Wallet Metrics
	walletModuleLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_module_loaded", Help: "Is the wallet module loaded. 0=not loaded.  1=loaded"})
	walletLocked = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_locked", Help: "Is the wallet locked. 0=not locked.  1=locked"})
	walletConfirmedSiacoinBalanceHastings = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_confirmed_siacoin_balance_hastings", Help: "Wallet confirmed Siacoin balance (Hastings)"})
	walletConfirmedSiacoinBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_confirmed_siacoin_balance", Help: "Wallet confirmed Siacoin balance (Siacoins)"})
	walletSiafundBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_siafund_balance", Help: "Wallet Siafund balance"})
	walletSiafundClaimBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_siafund_claim_balance", Help: "Wallet Siafund claim balance"})
	walletNumAddresses = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_num_addresses", Help: "Number of wallet addresses being tracked by Sia"})

	// Gateway Metrics
	gatewayModuleLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gateway_module_loaded", Help: "Is the gateway module loaded. 0=not loaded.  1=loaded"})
	gatewayNumPeers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gateway_num_peers", Help: "gateway number of peers"})
	gatewayRateLimitDownload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gateway_rate_limit_download", Help: "gateway download ratelimit (bytes-per-second)"})
	gatewayRateLimitUpload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gateway_rate_limit_upload", Help: "gateway upload ratelimit (bytes-per-second)"})

	// Hostdb Metrics
	hostdbNumAllHosts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostdb_num_all_hosts", Help: "Total number of hosts in hostdb"})
	hostdbNumActiveHosts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostdb_num_active_hosts", Help: "Number of active hosts in hostdb"})
	hostdbNumInactiveHosts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostdb_num_inactive_hosts", Help: "Number of inactive hosts in hostdb"})
	hostdbNumOfflineHosts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostdb_num_offline_hosts", Help: "Number of offline hosts in hostdb"})

	// Host Metrics
	hostAcceptingContracts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_accepting_contracts", Help: "Is the host accepting contracts 0=no, 1=yes"})
	hostMaxDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_max_duration", Help: "max duration in weeks"})
	hostMaxDownloadBatchSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_max_download_batch_size", Help: "Max Download Batch Size"})
	hostMaxReviseBatchSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_max_revise_batch_size", Help: "Max revise Batch Size"})
	hostWindowSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_window_size", Help: "Window Size in hours"})
	hostCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_collateral", Help: "Host Collateral in Siacoins"})
	hostCollateralBudget = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_collateral_budget", Help: "Host Collateral budget in Siacoins"})
	hostMaxCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_max_collateral", Help: "Max collateral per contract"})
	hostContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_contract_count", Help: "number of host contracts"})
	hostTotalStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_total_storage", Help: "total amount of storage available on the host in bytes"})
	hostRemainingStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_remaining_storage", Help: "amount of storage remaining on the host in bytes"})
)

const (
	moduleNotReadyStatus = "Module not loaded or still starting up"
)

func hostMetrics(sc *sia.Client) {
	hg, err := sc.HostGet()
	if errors.Contains(err, ErrAPICallNotRecognized) {
		// Assume module is not loaded if status command is not recognized.
		log.Info("Host module is not loaded")
		return
	} else if err != nil {
		log.Info("Could not fetch host settings")
	}

	sg, err := sc.HostStorageGet()
	if err != nil {
		log.Info("Could not fetch storage info")
	}

	es := hg.ExternalSettings
	fm := hg.FinancialMetrics
	is := hg.InternalSettings
	//	nm := hg.NetworkMetrics

	// calculate total storage available and remaining
	var totalstorage, storageremaining uint64
	for _, folder := range sg.Folders {
		totalstorage += folder.Capacity
		storageremaining += folder.CapacityRemaining
	}

	// convert price from bytes/block to TB/Month
	//	price := is.MinStoragePrice.Mul(modules.BlockBytesPerMonthTerabyte)
	// calculate total revenue
	//	totalRevenue := fm.ContractCompensation.
	//		Add(fm.StorageRevenue).
	//		Add(fm.DownloadBandwidthRevenue).
	//		Add(fm.UploadBandwidthRevenue)
	//	totalPotentialRevenue := fm.PotentialContractCompensation.
	//		Add(fm.PotentialStorageRevenue).
	//		Add(fm.PotentialDownloadBandwidthRevenue).
	//		Add(fm.PotentialUploadBandwidthRevenue)

	// Host Internal Settings
	hostAcceptingContracts.Set(boolToFloat64(is.AcceptingContracts))
	hostTotalStorage.Set(float64(es.TotalStorage))
	hostRemainingStorage.Set(float64(es.RemainingStorage))
	hostMaxDuration.Set(float64(is.MaxDuration))
	hostMaxDownloadBatchSize.Set(float64(is.MaxDownloadBatchSize))
	hostMaxReviseBatchSize.Set(float64(is.MaxReviseBatchSize))
	hostWindowSize.Set(float64(is.WindowSize / 6))
	hostCollateralFloat, _ := is.Collateral.Mul(modules.BlockBytesPerMonthTerabyte).Float64()
	hostCollateral.Set(hostCollateralFloat / 1e24)
	hostCollateralBudgetFloat, _ := is.CollateralBudget.Float64()
	hostCollateralBudget.Set(hostCollateralBudgetFloat / 1e24)
	hostMaxCollateralFloat, _ := is.MaxCollateral.Float64()
	hostMaxCollateral.Set(hostMaxCollateralFloat / 1e24)

	hostContractCount.Set(float64(fm.ContractCount))

}

func renterMetrics(sc *sia.Client) {

	// Renter Get Dir Metrics
	rg, err := sc.RenterGetDir(modules.RootSiaPath())
	if errors.Contains(err, ErrAPICallNotRecognized) {
		log.Info("Renter module is not loaded")
		renterModuleLoaded.Set(boolToFloat64(false))
		return
	} else if err != nil {
		log.Info("Could not get Renter metrics")
		return
	}

	renterModuleLoaded.Set(boolToFloat64(true))
	renterAggregateNumFiles.Set(float64(rg.Directories[0].AggregateNumFiles))
	renterAggregateNumStuckChunks.Set(float64(rg.Directories[0].AggregateNumStuckChunks))
	renterAggregateSize.Set(float64(rg.Directories[0].AggregateSize))
	renterMaxHealth.Set(float64(rg.Directories[0].MaxHealth))
	renterMinRedundancy.Set(float64(rg.Directories[0].MinRedundancy))

	// Contract Metrics
	rc, err := sc.RenterDisabledContractsGet()
	if err != nil {
		log.Info("Could not get renter contracts")
	}
	renterNumActiveContracts.Set(float64(len(rc.ActiveContracts)))
	renterNumPassiveContracts.Set(float64(len(rc.PassiveContracts)))
	renterNumRefreshedContracts.Set(float64(len(rc.RefreshedContracts)))
	renterNumDisabledContracts.Set(float64(len(rc.DisabledContracts)))

	rce, err := sc.RenterDisabledContractsGet()
	if err != nil {
		log.Info("Could not get renter expired contracts")
	}
	renterNumExpiredContracts.Set(float64(len(rce.ExpiredContracts)))
	renterNumExpiredRefreshedContracts.Set(float64(len(rce.ExpiredRefreshedContracts)))

	// Allowance Metrics
	ra, err := sc.RenterGet()
	if err != nil {
		log.Info("Could not get renter allowance info")
	}
	allowance := ra.Settings.Allowance
	funds, _ := allowance.Funds.Float64()
	renterAllowanceAmount.Set(float64(funds / 1e24))
	renterAllowancePeriod.Set(float64(allowance.Period))
	renterAllowanceRenewWindow.Set(float64(allowance.RenewWindow))
	renterAllowanceHosts.Set(float64(allowance.Hosts))

	fm := ra.FinancialMetrics
	totalSpent := fm.ContractFees.Add(fm.UploadSpending).Add(fm.DownloadSpending).Add(fm.StorageSpending)
	totalSpentFloat, _ := totalSpent.Float64()
	renterAllowanceCurrentSpent.Set(totalSpentFloat / 1e24)
	storageSpendingFloat, _ := fm.StorageSpending.Float64()
	renterAllowanceCurrentStorage.Set(storageSpendingFloat / 1e24)
	uploadSpendingFloat, _ := fm.UploadSpending.Float64()
	renterAllowanceCurrentUpload.Set(uploadSpendingFloat / 1e24)
	downloadSpendingFloat, _ := fm.DownloadSpending.Float64()
	renterAllowanceCurrentDownload.Set(downloadSpendingFloat / 1e24)
	contractFeesFloat, _ := fm.ContractFees.Float64()
	renterAllowanceCurrentFees.Set(contractFeesFloat / 1e24)
	unspentFloat, _ := fm.Unspent.Float64()
	renterAllowanceCurrentUnspent.Set(unspentFloat / 1e24)
	unspentAllocatedFloat, _ := fm.TotalAllocated.Sub(totalSpent).Float64()
	renterAllowanceCurrentUnspentAllocated.Set(unspentAllocatedFloat / 1e24)
	unspentUnallocatedFloat, _ := fm.Unspent.Sub(fm.TotalAllocated.Sub(totalSpent)).Float64()
	renterAllowanceCurrentUnspentUnallocated.Set(unspentUnallocatedFloat / 1e24)

	renterRateLimitUpload.Set(float64(ra.Settings.MaxUploadSpeed))
	renterRateLimitDownload.Set(float64(ra.Settings.MaxDownloadSpeed))

}

// consensuMetrics retrieves and sets the Prometheus metrics related to the
// consensus module
func consensusMetrics(sc *sia.Client) {
	cs, err := sc.ConsensusGet()
	if errors.Contains(err, ErrAPICallNotRecognized) {
		log.Info("Consensus module is not loaded")
		consensusModuleLoaded.Set(boolToFloat64(false))
		return
	} else if err != nil {
		log.Info("Could not get Consensus metrics")
		return
	}

	consensusModuleLoaded.Set(boolToFloat64(true))
	consensusSynced.Set(boolToFloat64(cs.Synced))
	consensusHeight.Set(float64(cs.Height))
	Difficulty, _ := cs.Difficulty.Float64()
	consensusDifficulty.Set(Difficulty)
	consensusHshrt.Set(Difficulty/600.0)
}

// daemonMetrics retrieves and sets the Prometheus metrics related to the
// Sia daemon
func daemonMetrics(sc *sia.Client) {
	//al, err := sc.DaemonAlertsGet()
	//if err != nil {
	//	log.Info("Could not get Daemon metrics")
	//	return
	//}
	//daemonAggregateNumAlerts.Set(float64(len(al.Alerts)))

	// Global Daemon Rate Limits
	dg, err := sc.DaemonSettingsGet()
	if err != nil {
		log.Info("Could not get daemon metrics")
		return
	}
	daemonRateLimitUpload.Set(float64(dg.MaxUploadSpeed))
	daemonRateLimitDownload.Set(float64(dg.MaxDownloadSpeed))

}

// walletMetrics retrieves and sets the Prometheus metrics related to the
// Sia wallet
func walletMetrics(sc *sia.Client) {
	status, err := sc.WalletGet()
	if errors.Contains(err, ErrAPICallNotRecognized) {
		log.Info("Wallet module is not loaded")
		walletModuleLoaded.Set(boolToFloat64(false))
		return
	} else if err != nil {
		log.Info("Could not get Wallet metrics")
		return
	}
	walletModuleLoaded.Set(boolToFloat64(true))
	if !status.Unlocked {
		walletLocked.Set(boolToFloat64(false))
	}
	walletLocked.Set(boolToFloat64(true))

	ConfirmedBalance, _ := status.ConfirmedSiacoinBalance.Float64()
	walletConfirmedSiacoinBalanceHastings.Set(ConfirmedBalance)
	walletConfirmedSiacoinBalance.Set(ConfirmedBalance / 1e24)

	SiafundBalance, _ := status.SiafundBalance.Float64()
	walletSiafundBalance.Set(SiafundBalance)

	SiafundClaimBalance, _ := status.SiacoinClaimBalance.Float64()
	walletSiafundClaimBalance.Set(SiafundClaimBalance)

	addresses, err := sc.WalletAddressesGet()
	if err != nil {
		log.Info("Could not get wallet addresses")
	}
	walletNumAddresses.Set(float64(len(addresses.Addresses)))
}

// gatewayMetrics retrieves and sets the Prometheus metrics related to the
// Sia gateway
func gatewayMetrics(sc *sia.Client) {
	gateway, err := sc.GatewayGet()
	if errors.Contains(err, ErrAPICallNotRecognized) {
		log.Info("Gateway module is not loaded")
		gatewayModuleLoaded.Set(boolToFloat64(false))
		return
	} else if err != nil {
		log.Info("Could not get Gateway metrics")
		return
	}

	gatewayModuleLoaded.Set(boolToFloat64(true))
	gatewayNumPeers.Set(float64(len(gateway.Peers)))
	gatewayRateLimitUpload.Set(float64(gateway.MaxUploadSpeed))
	gatewayRateLimitDownload.Set(float64(gateway.MaxDownloadSpeed))
}

// hostdbMetrics retrieves and sets the Prometheus metrics related to the
// Sia hostdb
func hostdbMetrics(sc *sia.Client) {
	hostdb, err := sc.HostDbAllGet()
	if errors.Contains(err, ErrAPICallNotRecognized) {
		log.Info("HostDB module is not loaded")
		return
	} else if err != nil {
		log.Info("Could not get Gateway metrics")
		return
	}

	// Iterate through the hosts and divide by category.
	var activeHosts, inactiveHosts, offlineHosts []api.ExtendedHostDBEntry
	for _, host := range hostdb.Hosts {
		if host.AcceptingContracts && len(host.ScanHistory) > 0 && host.ScanHistory[len(host.ScanHistory)-1].Success {
			activeHosts = append(activeHosts, host)
			continue
		}
		if len(host.ScanHistory) > 0 && host.ScanHistory[len(host.ScanHistory)-1].Success {
			inactiveHosts = append(inactiveHosts, host)
			continue
		}
		offlineHosts = append(offlineHosts, host)
	}

	hostdbNumAllHosts.Set(float64(len(hostdb.Hosts)))
	hostdbNumActiveHosts.Set(float64(len(activeHosts)))
	hostdbNumInactiveHosts.Set(float64(len(inactiveHosts)))
	hostdbNumOfflineHosts.Set(float64(len(offlineHosts)))
}
